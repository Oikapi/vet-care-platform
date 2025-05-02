from flask import Blueprint, request, jsonify
from services.pet_service import register_pet, find_pet_by_id

pet_bp = Blueprint('pet', __name__)

@pet_bp.route('/register', methods=['POST'])
def register():
    """
    Маршрут для регистрации нового питомца.
    """
    data = request.json
    try:
        pet_id = register_pet(
            species=data.get("species"),
            name=data.get("name"),
            breed=data.get("breed"),
            gender=data.get("gender"),
            age=data.get("age")
        )
        return jsonify({"message": "Pet registered", "id": pet_id}), 201
    except Exception as e:
        return jsonify({"error": str(e)}), 500

@pet_bp.route('/find_pet/<int:pet_id>', methods=['GET'])
def find_pet(pet_id):
    """
    Маршрут для поиска питомца по ID.
    """
    pet = find_pet_by_id(pet_id)
    if pet:
        return jsonify(pet), 200
    else:
        return jsonify({"error": "Pet not found"}), 404