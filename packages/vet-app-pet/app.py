from flask import Flask
from routes.pet_routes import pet_bp  

app = Flask(__name__) 

app.register_blueprint(pet_bp, url_prefix='')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)