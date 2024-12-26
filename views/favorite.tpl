<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Favorite Cat Images</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
            margin: 0;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }

        .favorite-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }

        .favorite-card {
            background: white;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }

        .favorite-card img {
            width: 100%;
            height: 300px;
            object-fit: cover;
        }

        .back-button {
            display: inline-block;
            padding: 10px 20px;
            background-color: #FF4D4D;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            margin-bottom: 20px;
        }

        .no-favorites {
            text-align: center;
            padding: 50px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/cat" class="back-button">
            <i class="fas fa-arrow-left"></i> Back to Gallery
        </a>
        
        <div class="favorite-grid">
            {{if .CatFavourites}}
                {{range .CatFavourites}}
                    <div class="favorite-card">
                        <img src="{{.Image.URL}}" alt="Favorite Cat">
                    </div>
                {{end}}
            {{else}}
                <div class="no-favorites">
                    <h2>No favorites yet</h2>
                    <p>Start adding some favorite cats!</p>
                </div>
            {{end}}
        </div>
    </div>
</body>
</html>