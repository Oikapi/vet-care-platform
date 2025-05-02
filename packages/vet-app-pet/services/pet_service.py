from db.db_connection import get_db_connection

def register_pet(species, name, breed, gender, age):
    """
    Регистрация нового питомца.
    """
    connection = get_db_connection()
    cursor = connection.cursor()

    query = """
    INSERT INTO pets (species, name, breed, gender, age)
    VALUES (%s, %s, %s, %s, %s)
    """
    values = (species, name, breed, gender, age)
    cursor.execute(query, values)
    connection.commit()

    pet_id = cursor.lastrowid
    cursor.close()
    connection.close()
    return pet_id

def find_pet_by_id(pet_id):
    """
    Поиск питомца по ID.
    """
    connection = get_db_connection()
    cursor = connection.cursor(dictionary=True)

    query = "SELECT * FROM pets WHERE id = %s"
    cursor.execute(query, (pet_id,))
    pet = cursor.fetchone()

    cursor.close()
    connection.close()
    return pet