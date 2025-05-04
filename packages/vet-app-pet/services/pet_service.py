from db.db_connection import get_db_connection
import mysql.connector


def register_pet(species, name, breed, gender, age, user_id):
    """Регистрация нового питомца"""
    conn = None
    cursor = None

    try:
        conn = get_db_connection()
        cursor = conn.cursor()

        query = """
        INSERT INTO pets (species, name, breed, gender, age, user_id)
        VALUES (%s, %s, %s, %s, %s, %s)
        """
        values = (species, name, breed, gender, age, user_id)
        cursor.execute(query, values)
        conn.commit()

        return cursor.lastrowid
    except mysql.connector.Error as e:
        if conn:
            conn.rollback()
        raise Exception(f"Database error: {e}")
    finally:
        if cursor:
            cursor.close()
        if conn:
            conn.close()


def find_pet_by_id(pet_id):
    """Получить питомца по ID"""
    conn = None
    cursor = None

    try:
        conn = get_db_connection()
        cursor = conn.cursor(dictionary=True)

        query = "SELECT * FROM pets WHERE id = %s"
        cursor.execute(query, (pet_id,))
        return cursor.fetchone()
    except mysql.connector.Error as e:
        raise Exception(f"Database error: {e}")
    finally:
        if cursor:
            cursor.close()
        if conn:
            conn.close()


def get_all_pets():
    """Получить всех питомцев"""
    conn = None
    cursor = None

    try:
        conn = get_db_connection()
        cursor = conn.cursor(dictionary=True)

        query = "SELECT * FROM pets"
        cursor.execute(query)
        return cursor.fetchall()
    except mysql.connector.Error as e:
        raise Exception(f"Database error: {e}")
    finally:
        if cursor:
            cursor.close()
        if conn:
            conn.close()


def update_pet(pet_id, data):
    """Обновление информации о питомце"""
    conn = None
    cursor = None

    try:
        conn = get_db_connection()
        cursor = conn.cursor()

        allowed_fields = ["species", "name", "breed", "gender", "age", "user_id"]
        updates = []
        values = []

        for field in allowed_fields:
            if field in data:
                updates.append(f"{field} = %s")
                values.append(data[field])

        if not updates:
            raise ValueError("No valid fields to update")

        values.append(pet_id)
        query = f"UPDATE pets SET {', '.join(updates)} WHERE id = %s"
        cursor.execute(query, values)
        conn.commit()

        return cursor.rowcount > 0
    except mysql.connector.Error as e:
        if conn:
            conn.rollback()
        raise Exception(f"Database error: {e}")
    finally:
        if cursor:
            cursor.close()
        if conn:
            conn.close()