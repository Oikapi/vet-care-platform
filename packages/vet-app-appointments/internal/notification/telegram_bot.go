package notification

import (
	"fmt"
	"vet-app-appointments/internal/models"
	"vet-app-appointments/internal/service"
	"vet-app-appointments/pkg/logger"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

type TelegramBot struct {
	bot        *tgbotapi.BotAPI
	service    *service.AppointmentService
	log        *logger.Logger
	userStates map[int64]*BookingState
}

type BookingState struct {
	Step     int
	ClientID uint
	DoctorID uint
	ClinicID uint
	SlotID   uint
}

func NewTelegramBot(token string, service *service.AppointmentService) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Telegram бота: %v", err)
	}
	return &TelegramBot{
		bot:        bot,
		service:    service,
		log:        logger.NewLogger(),
		userStates: make(map[int64]*BookingState),
	}, nil
}

// Экранирование специальных символов для MarkdownV2
func escapeMarkdownV2(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

func (t *TelegramBot) Start() {
	t.log.Info("Запуск Telegram бота...")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			t.handleCallbackQuery(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chatID, "")
		msg.ParseMode = "MarkdownV2"

		if state, exists := t.userStates[chatID]; exists {
			t.handleBookingSteps(chatID, update.Message.Text, state)
			continue
		}

		if update.Message.IsCommand() {
			command := update.Message.Command()
			args := update.Message.CommandArguments()

			switch command {
			case "start":
				msg.Text = "🐾 *Привет\\!* Я бот, который поможет тебе записаться к ветеринару\\! 🎉\n\n" +
					"Вот что я умею:\n" +
					"📅 *Посмотреть доступные слоты*: `/slots <id_клиники> <дата>`\n" +
					"   Пример: `/slots 1 2025\\-05\\-02`\n" +
					"📋 *Узнать о записи*: `/appointment <id_записи>`\n" +
					"   Пример: `/appointment 1`\n" +
					"✍️ *Создать запись*: `/book <id_клиента> <id_доктора> <id_клиники> <id_слота>`\n" +
					"   Пример: `/book 1 1 1 1`\n" +
					"🚫 *Отменить запись*: `/cancel <id_записи>`\n" +
					"   Пример: `/cancel 1`\n\n" +
					"Выбери команду ниже или напиши её вручную\\! 😊"
				msg.ReplyMarkup = createMainMenu()
			case "slots":
				t.handleSlotsCommand(chatID, args)
				continue
			case "appointment":
				t.handleAppointmentCommand(chatID, args)
				continue
			case "book":
				t.handleBookCommand(chatID, args)
				continue
			case "cancel":
				t.handleCancelCommand(chatID, args)
				continue
			default:
				msg.Text = "🤔 Неизвестная команда\\. Напиши `/start`, чтобы увидеть список доступных команд\\!"
			}
		} else {
			switch update.Message.Text {
			case "📅 Посмотреть слоты":
				msg.Text = "📅 Введи команду в формате: `/slots <id_клиники> <дата>`\nПример: `/slots 1 2025\\-05\\-02`"
			case "📋 Мои записи":
				msg.Text = "📋 Введи команду в формате: `/appointment <id_записи>`\nПример: `/appointment 1`"
			case "✍️ Создать запись":
				t.userStates[chatID] = &BookingState{Step: 1}
				msg.Text = "✍️ Давай создадим запись\\! Шаг 1: Введи ID клиента \\(например, 1\\)"
			case "🚫 Отменить запись":
				msg.Text = "🚫 Введи команду в формате: `/cancel <id_записи>`\nПример: `/cancel 1`"
			default:
				msg.Text = "🤔 Не понимаю\\. Используй кнопки ниже или напиши `/start` для списка команд\\!"
				msg.ReplyMarkup = createMainMenu()
			}
		}

		if _, err := t.bot.Send(msg); err != nil {
			t.log.Errorf("Ошибка отправки сообщения: %v", err)
		}
	}
}

func (t *TelegramBot) handleBookingSteps(chatID int64, input string, state *BookingState) {
	msg := tgbotapi.NewMessage(chatID, "")
	msg.ParseMode = "MarkdownV2"

	switch state.Step {
	case 1:
		clientID, err := strconv.Atoi(input)
		if err != nil {
			msg.Text = "❌ Ошибка: ID клиента должен быть числом\\. Попробуй снова\\!"
		} else {
			state.ClientID = uint(clientID)
			state.Step = 2
			msg.Text = "✍️ Шаг 2: Введи ID доктора \\(например, 1\\)"
		}
	case 2:
		doctorID, err := strconv.Atoi(input)
		if err != nil {
			msg.Text = "❌ Ошибка: ID доктора должен быть числом\\. Попробуй снова\\!"
		} else {
			state.DoctorID = uint(doctorID)
			state.Step = 3
			msg.Text = "✍️ Шаг 3: Введи ID клиники \\(например, 1\\)"
		}
	case 3:
		clinicID, err := strconv.Atoi(input)
		if err != nil {
			msg.Text = "❌ Ошибка: ID клиники должен быть числом\\. Попробуй снова\\!"
		} else {
			state.ClinicID = uint(clinicID)
			state.Step = 4
			msg.Text = "✍️ Шаг 4: Введи ID слота \\(например, 1\\)"
		}
	case 4:
		slotID, err := strconv.Atoi(input)
		if err != nil {
			msg.Text = "❌ Ошибка: ID слота должен быть числом\\. Попробуй снова\\!"
		} else {
			state.SlotID = uint(slotID)
			appointment := &models.Appointment{
				ClientID:   state.ClientID,
				DoctorID:   state.DoctorID,
				ClinicID:   state.ClinicID,
				SlotID:     state.SlotID,
				TelegramID: strconv.FormatInt(chatID, 10),
			}

			if err := t.service.CreateAppointment(appointment); err != nil {
				msg.Text = sensibleError(fmt.Sprintf("❌ Ошибка создания записи: %v", err))
			} else {
				msg.Text = fmt.Sprintf("🎉 Запись успешно создана\\! ID: %d", appointment.ID)
			}
			delete(t.userStates, chatID)
			msg.ReplyMarkup = createMainMenu()
		}
	}

	if _, err := t.bot.Send(msg); err != nil {
		t.log.Errorf("Ошибка отправки сообщения: %v", err)
	}
}

func createMainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📅 Посмотреть слоты"),
			tgbotapi.NewKeyboardButton("📋 Мои записи"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✍️ Создать запись"),
			tgbotapi.NewKeyboardButton("🚫 Отменить запись"),
		),
	)
}

func (t *TelegramBot) handleSlotsCommand(chatID int64, args string) {
	parts := strings.Fields(args)
	if len(parts) != 2 {
		t.sendMessage(chatID, "📅 Использование: `/slots <id_клиники> <дата>`\nПример: `/slots 1 2025\\-05\\-02`")
		return
	}

	clinicID, err := strconv.Atoi(parts[0])
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: id_клиники должен быть числом")
		return
	}

	date, err := time.Parse("2006-01-02", parts[1])
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: дата должна быть в формате ГГГГ-ММ-ДД \\(например, 2025\\-05\\-02\\)")
		return
	}

	slots, err := t.service.GetAvailableSlots(uint(clinicID), date)
	if err != nil {
		t.sendMessage(chatID, fmt.Sprintf("❌ Ошибка получения слотов: %v", err))
		return
	}

	if len(slots) == 0 {
		t.sendMessage(chatID, "😔 Доступных слотов нет\\.")
		return
	}

	var response strings.Builder
	response.WriteString("*Доступные слоты:*\n")

	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	for _, slot := range slots {
		slotInfo := fmt.Sprintf("ID: %d, Время: %s", slot.ID, slot.SlotTime.Format("2006-01-02 15:04"))
		response.WriteString(escapeMarkdownV2(slotInfo) + "\n")

		callbackData := fmt.Sprintf("book_slot:%d:%d", clinicID, slot.ID)
		button := tgbotapi.NewInlineKeyboardButtonData(slotInfo, callbackData)
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(button))
	}

	msg := tgbotapi.NewMessage(chatID, response.String())
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
	if _, err := t.bot.Send(msg); err != nil {
		t.log.Errorf("Ошибка отправки сообщения: %v", err)
	}
}

func (t *TelegramBot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := strings.Split(callback.Data, ":")

	if len(data) != 3 || data[0] != "book_slot" {
		t.sendMessage(chatID, "❌ Ошибка: неверный формат команды")
		return
	}

	clinicID, err := strconv.Atoi(data[1])
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: clinic_id должен быть числом")
		return
	}

	slotID, err := strconv.Atoi(data[2])
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: slot_id должен быть числом")
		return
	}

	appointment := &models.Appointment{
		ClientID:   1,
		DoctorID:   1,
		ClinicID:   uint(clinicID),
		SlotID:     uint(slotID),
		TelegramID: strconv.FormatInt(chatID, 10),
	}

	if err := t.service.CreateAppointment(appointment); err != nil {
		t.sendMessage(chatID, sensibleError(fmt.Sprintf("❌ Ошибка создания записи: %v", err)))
		return
	}

	t.sendMessage(chatID, fmt.Sprintf("🎉 Запись успешно создана\\! ID: %d", appointment.ID))

	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := t.bot.Request(callbackConfig); err != nil {
		t.log.Errorf("Ошибка подтверждения callback: %v", err)
	}
}

func (t *TelegramBot) handleAppointmentCommand(chatID int64, args string) {
	appointmentID, err := strconv.Atoi(args)
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: ID записи должен быть числом")
		return
	}

	appointment, err := t.service.GetAppointment(uint(appointmentID))
	if err != nil {
		t.sendMessage(chatID, sensibleError(fmt.Sprintf("❌ Ошибка получения записи: %v", err)))
		return
	}

	if appointment.TelegramID != strconv.FormatInt(chatID, 10) {
		t.sendMessage(chatID, "🔒 Ошибка: у вас нет доступа к этой записи")
		return
	}

	response := fmt.Sprintf("*Запись #%d:*\n👤 Клиент: %d\n👨‍⚕️ Доктор: %d\n🏥 Клиника: %d\n⏰ Слот: %d \\(%s\\)\n📌 Статус: %s",
		appointment.ID, appointment.ClientID, appointment.DoctorID, appointment.ClinicID,
		appointment.SlotID, escapeMarkdownV2(appointment.Slot.SlotTime.Format("2006-01-02 15:04")), escapeMarkdownV2(appointment.Status))
	t.sendMessageWithMarkdown(chatID, response)
}

func (t *TelegramBot) handleBookCommand(chatID int64, args string) {
	parts := strings.Fields(args)
	if len(parts) != 4 {
		t.sendMessage(chatID, "✍️ Использование: `/book <id_клиента> <id_доктора> <id_клиники> <id_слота>`\nПример: `/book 1 1 1 1`")
		return
	}

	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: id_клиента должен быть числом")
		return
	}
	doctorID, err := strconv.Atoi(parts[1])
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: id_доктора должен быть числом")
		return
	}
	clinicID, err := strconv.Atoi(parts[2])
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: id_клиники должен быть числом")
		return
	}
	slotID, err := strconv.Atoi(parts[3])
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: id_слота должен быть числом")
		return
	}

	appointment := &models.Appointment{
		ClientID:   uint(clientID),
		DoctorID:   uint(doctorID),
		ClinicID:   uint(clinicID),
		SlotID:     uint(slotID),
		TelegramID: strconv.FormatInt(chatID, 10),
	}

	if err := t.service.CreateAppointment(appointment); err != nil {
		t.sendMessage(chatID, sensibleError(fmt.Sprintf("❌ Ошибка создания записи: %v", err)))
		return
	}

	t.sendMessage(chatID, fmt.Sprintf("🎉 Запись успешно создана\\! ID: %d", appointment.ID))
}

func (t *TelegramBot) handleCancelCommand(chatID int64, args string) {
	appointmentID, err := strconv.Atoi(args)
	if err != nil {
		t.sendMessage(chatID, "❌ Ошибка: ID записи должен быть числом")
		return
	}

	appointment, err := t.service.GetAppointment(uint(appointmentID))
	if err != nil {
		t.sendMessage(chatID, sensibleError(fmt.Sprintf("❌ Ошибка получения записи: %v", err)))
		return
	}

	if appointment.TelegramID != strconv.FormatInt(chatID, 10) {
		t.sendMessage(chatID, "🔒 Ошибка: у вас нет доступа к этой записи")
		return
	}

	appointment.Status = "cancelled"
	if err := t.service.UpdateAppointment(appointment); err != nil {
		t.sendMessage(chatID, fmt.Sprintf("❌ Ошибка отмены записи: %v", err))
		return
	}

	if err := t.service.UnbookSlot(appointment.SlotID); err != nil {
		t.sendMessage(chatID, fmt.Sprintf("❌ Ошибка освобождения слота: %v", err))
		return
	}

	t.sendMessage(chatID, fmt.Sprintf("✅ Запись #%d успешно отменена", appointment.ID))
}

func (t *TelegramBot) sendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := t.bot.Send(msg); err != nil {
		t.log.Errorf("Ошибка отправки сообщения: %v", err)
	}
}

func (t *TelegramBot) sendMessageWithMarkdown(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "MarkdownV2"
	if _, err := t.bot.Send(msg); err != nil {
		t.log.Errorf("Ошибка отправки сообщения: %v", err)
	}
}

func (t *TelegramBot) SendReminder(chatID string, message string) error {
	id, err := chatIDToInt64(chatID)
	if err != nil {
		return fmt.Errorf("ошибка преобразования chatID: %v", err)
	}
	t.sendMessage(id, message)
	return nil
}

func chatIDToInt64(chatID string) (int64, error) {
	id, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат chatID: %v", err)
	}
	return id, nil
}

func sensibleError(err string) string {
	if strings.Contains(err, "record not found") {
		return "Запись не найдена"
	}
	return escapeMarkdownV2(err)
}
