from flask import Blueprint, request, jsonify
from services.pet_service import register_pet, find_pet_by_id, get_all_pets, update_pet

pet_bp = Blueprint('pet', __name__)

# --- Существующие маршруты ---
@pet_bp.route('/register', methods=['POST'])
def register():
    data = request.json
    try:
        pet_id = register_pet(
            species=data.get("species"),
            name=data.get("name"),
            breed=data.get("breed"),
            gender=data.get("gender"),
            age=data.get("age"),
            user_id=data.get("user_id"),
        )
        return jsonify({"message": "Pet registered", "id": pet_id}), 201
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@pet_bp.route('/find_pet/<int:pet_id>', methods=['GET'])
def find_pet(pet_id):
    pet = find_pet_by_id(pet_id)
    if pet:
        return jsonify(pet), 200
    else:
        return jsonify({"error": "Pet not found"}), 404


@pet_bp.route('/update/<int:pet_id>', methods=['PUT'])
def update_pet_route(pet_id):
    data = request.json
    if not data:
        return jsonify({"error": "No data provided"}), 400

    try:
        success = update_pet(pet_id, data)
        if success:
            return jsonify({"message": "Pet updated successfully"}), 200
        return jsonify({"error": "Pet not found or no changes made"}), 404
    except ValueError as e:
        return jsonify({"error": str(e)}), 400
    except Exception as e:
        return jsonify({"error": str(e)}), 500


# --- Новый маршрут: получить всех питомцев ---
@pet_bp.route('/pets', methods=['GET'])
def list_all_pets():
    try:
        pets = get_all_pets()
        if not pets:
            return jsonify({"message": "No pets found"}), 404
        return jsonify(pets), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500