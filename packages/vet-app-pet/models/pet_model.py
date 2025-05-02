class Pet:
    def __init__(self, species, age, breed):
        self.species = species
        self.age = age
        self.breed = breed

    def to_dict(self):
        """
        Преобразование объекта в словарь для работы с базой данных.
        """
        return {
            "species": self.species,
            "age": self.age,
            "breed": self.breed,
        }