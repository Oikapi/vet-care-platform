from flask import Flask, request, jsonify
import mysql.connector
from mysql.connector import Error
import os
app = Flask(__name__)

# Подключение к базе данных
def get_db_connection():
    try:
        connection = mysql.connector.connect(
            host=os.environ.get("DB_HOST", "localhost"),
            user=os.environ.get("DB_USER", "user"),
            password=os.environ.get("DB_PASSWORD", "password"),
            database=os.environ.get("DB_NAME", "pet_profiles")
        )
        if connection.is_connected():
            print("Successfully connected to the database")
        return connection
    except Error as err:
        print(f"Database connection error: {err}")
        raise

# Маршрут для добавления питомца
@app.route('/add_pet', methods=['POST'])
def add_pet():
    data = request.json

    required_fields = ["species", "name", "breed", "gender", "age", "user_id"]
    for field in required_fields:
        if field not in data:
            return jsonify({"error": f"Missing field: {field}"}), 400

    connection = get_db_connection()
    cursor = connection.cursor()

    try:
        query = """
        INSERT INTO pets (species, name, breed, gender, age, user_id)
        VALUES (%s, %s, %s, %s, %s, %s)
        """
        values = (
            data.get("species"),
            data.get("name"),
            data.get("breed"),
            data.get("gender"),
            data.get("age"),
            data.get("user_id")
        )
        cursor.execute(query, values)
        connection.commit()

        pet_id = cursor.lastrowid
        return jsonify({"message": "Pet added", "id": pet_id}), 201

    except mysql.connector.Error as e:
        connection.rollback()
        return jsonify({"error": str(e)}), 500

    finally:
        cursor.close()
        connection.close()

# Маршрут для получения питомца с указанием владельца
@app.route('/find_pet/<int:pet_id>', methods=['GET'])
def find_pet(pet_id):
    connection = get_db_connection()
    cursor = connection.cursor(dictionary=True)

    try:
        query = "SELECT * FROM pets WHERE id = %s"
        cursor.execute(query, (pet_id,))
        pet = cursor.fetchone()

        if not pet:
            return jsonify({"error": "Pet not found"}), 404

        # Добавляем информацию о владельце (если есть API пользователей)
        user_id = pet["user_id"]
        pet["owner"] = {"user_id": user_id}  # Простой пример, можно расширить через API

        return jsonify(pet), 200

    finally:
        cursor.close()
        connection.close()

# Новый маршрут для редактирования питомца
@app.route('/update_pet/<int:pet_id>', methods=['PUT'])
def update_pet(pet_id):
    data = request.json
    
    # Проверяем, есть ли данные для обновления
    if not data:
        return jsonify({"error": "No data provided"}), 400

    connection = get_db_connection()
    cursor = connection.cursor()

    try:
        # Формируем запрос динамически на основе переданных полей
        updates = []
        values = []
        allowed_fields = ["species", "name", "breed", "gender", "age", "user_id"]
        
        for field in allowed_fields:
            if field in data:
                updates.append(f"{field} = %s")
                values.append(data[field])
        
        if not updates:
            return jsonify({"error": "No valid fields to update"}), 400
        
        values.append(pet_id)  # Добавляем pet_id в конец для WHERE
        
        query = f"UPDATE pets SET {', '.join(updates)} WHERE id = %s"
        cursor.execute(query, values)
        connection.commit()

        if cursor.rowcount == 0:
            return jsonify({"error": "Pet not found or no changes made"}), 404

        return jsonify({"message": "Pet updated successfully"}), 200

    except mysql.connector.Error as e:
        connection.rollback()
        return jsonify({"error": str(e)}), 500

    finally:
        cursor.close()
        connection.close()

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)