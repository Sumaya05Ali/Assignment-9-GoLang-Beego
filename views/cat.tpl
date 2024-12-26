<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Images and Breeds</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
            margin: 0;
            padding: 20px;
            min-height: 100vh;
        }

        .container {
            background-color: white;
            border-radius: 15px;
            box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            margin: 160px auto;
            overflow: hidden;
        }

        .nav-buttons {
            display: flex;
            padding: 15px;
            border-bottom: 1px solid #eee;
        }

        .nav-button {
            display: flex;
            align-items: center;
            padding: 8px 15px;
            margin-right: 15px;
            border: none;
            background: none;
            cursor: pointer;
            color: #666;
            font-size: 14px;
        }

        .nav-button.active {
            color: #FF4D4D;
        }

        .nav-button i {
            margin-right: 8px;
        }

        .cat-image {
            position: relative;
            padding: 15px;
        }

        .cat-image img {
            width: 100%;
            height: 400px;
            object-fit: cover;
            border-radius: 10px;
        }

        .interaction-container {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 15px;
        }

        .favorite-button {
            background: none;
            border: none;
            cursor: pointer;
            padding: 8px;
            font-size: 20px;
            color: #666;
            margin-left: 15px;
        }
       
       .favorite-button.active {
            color: #FF4D4D;
       }
        .voting-buttons {
            display: flex;
            gap: 10px;
            margin-right: 15px;
        }

        .interaction-button {
            background: none;
            border: none;
            cursor: pointer;
            padding: 8px;
            font-size: 20px;
            color: #666;
        }

        .interaction-button:hover, .favorite-button:hover {
            color: #FF4D4D;
        }

        .loader {
            text-align: center;
            padding: 20px;
            color: #666;
        }
         .notification {
            position: fixed;
            top: 20px;
            right: 20px;
            background: #FF4D4D;
            color: white;
            padding: 15px 25px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.2);
            opacity: 0;
            transform: translateY(-20px);
            transition: all 0.3s ease;
        }

        .notification.show {
            opacity: 1;
            transform: translateY(0);
        }
        #breedsContent {
            display: none;
        }

        #votingContent {
            display: block;
        }
    </style>
</head>
<body>
     <div class="container">
        <div class="nav-buttons">
            <button class="nav-button" id="votingBtn">
                <i class="fas fa-arrow-up"></i>
                Voting
            </button>
            <button class="nav-button" id="breedBtn">
                <i class="fas fa-search"></i>
                Breeds
            </button>
            <button class="nav-button" id="favBtn">
                <i class="fas fa-heart"></i>
                Favs
            </button>
        </div>
    
        <div class="cat-image" id="catImageContainer">
            <div class="loader">Loading cat image...</div>
        </div>

        <div class="interaction-container">
            <div class="favorite-section">
                <button class="favorite-button" id="favoriteBtn">
                    <i class="far fa-heart"></i>
                </button>
            </div>
            <div class="voting-buttons">
                <button class="interaction-button" id="likeBtn">
                    <i class="far fa-thumbs-up"></i>
                </button>
                <button class="interaction-button" id="dislikeBtn">
                    <i class="far fa-thumbs-down"></i>
                </button>
            </div>
        </div>
    </div>
        <div id="breedsContent">
            <div id="breedSearchContainer"></div>
        </div>
    </div>
    <div id="notification" class="notification"></div>
    <script>
        const API_KEY = 'live_StuW3GXvHysXlRahhSSy22vXbKQAZUl6yYJxVhXnlpJ18BhC4Pinz9Sx1Xr7F6BT';
        const apiUrl = 'https://api.thecatapi.com/v1/images/search?limit=1';
        let currentImageId = null;

        // Function to fetch cat image
        async function fetchCatImage() {
            try {
                const response = await fetch(apiUrl, {
                    headers: {
                        'x-api-key': API_KEY
                    }
                });
                const [catImage] = await response.json();
                currentImageId = catImage.id;
                
                const imageContainer = document.getElementById('catImageContainer');
                imageContainer.innerHTML = `<img src="${catImage.url}" alt="Cat Image">`;
                
                // Reset favorite button state
                const favoriteBtn = document.getElementById('favoriteBtn');
                favoriteBtn.classList.remove('active');
                favoriteBtn.querySelector('i').className = 'far fa-heart';
            } catch (error) {
                console.error('Error fetching cat image:', error);
                document.getElementById('catImageContainer').innerHTML = 
                    '<p>Failed to load cat image. Please try again later.</p>';
            }
        }
    
     function showNotification(message) {
            const notification = document.getElementById('notification');
            notification.textContent = message;
            notification.classList.add('show');
            
            setTimeout(() => {
                notification.classList.remove('show');
            }, 3000);
        } 

       async function saveFavorite() {
            if (!currentImageId) return;

            try {
                const response = await fetch('/save-favorite', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        image_id: currentImageId
                    })
                });


               if (response.ok) {
                    const favoriteBtn = document.getElementById('favoriteBtn');
                    favoriteBtn.classList.add('active');
                    favoriteBtn.querySelector('i').className = 'fas fa-heart';
                    showNotification('Image saved to favorites!');
                    
                    // Fetch new image after successful favorite
                    setTimeout(() => {
                        fetchCatImage();
                    }, 300);
                } else {
                    throw new Error('Failed to save favorite');
                }
            } catch (error) {
                showNotification('Failed to save favorite. Please try again.');
            }
        }

        // Updated fetchCatImage function with logging
        async function fetchCatImage() {
            try {
                const response = await fetch(apiUrl, {
                    headers: {
                        'x-api-key': API_KEY
                    }
                });
                const [catImage] = await response.json();
                currentImageId = catImage.id;
                console.log('Fetched new image with ID:', currentImageId);
                
                const imageContainer = document.getElementById('catImageContainer');
                imageContainer.innerHTML = `<img src="${catImage.url}" alt="Cat Image" data-id="${catImage.id}">`;
                
                const favoriteBtn = document.getElementById('favoriteBtn');
                favoriteBtn.classList.remove('active');
                favoriteBtn.querySelector('i').className = 'far fa-heart';
            } catch (error) {
                console.error('Error fetching cat image:', error);
                document.getElementById('catImageContainer').innerHTML = 
                    '<p>Failed to load cat image. Please try again later.</p>';
            }
        }

        // Load cat image on page load
        window.onload = function() {
            fetchCatImage();
            handleCurrentPath();
        };

        // Voting button functionality
        document.getElementById('likeBtn').addEventListener('click', () => {
            console.log('Liked the image');
            fetchCatImage();
        });

        document.getElementById('dislikeBtn').addEventListener('click', () => {
            console.log('Disliked the image');
            fetchCatImage();
        });

        // Add favorite button functionality
        document.getElementById('favoriteBtn').addEventListener('click', saveFavorite);

        // Navigation functionality
        function handleCurrentPath() {
            const path = window.location.pathname.split('/').pop();
            setActiveButton(path || 'voting');
        }

        function setActiveButton(path) {
    document.querySelectorAll('.nav-button').forEach(btn => btn.classList.remove('active'));

    // Get elements and check if they exist
    const votingContent = document.getElementById('votingContent');
    const breedsContent = document.getElementById('breedsContent');

    switch (path) {
        case 'voting':
            document.getElementById('votingBtn').classList.add('active');

            if (votingContent) {
                votingContent.style.display = 'block';
            }
            if (breedsContent) {
                breedsContent.style.display = 'none';
            }
            break;

        case 'breeds':
            document.getElementById('breedBtn').classList.add('active');

            if (votingContent) {
                votingContent.style.display = 'none';
            }
            if (breedsContent) {
                breedsContent.style.display = 'block';
                loadBreedsContent(); // Make sure this function exists and works
            }
            break;

        case 'favs':
            document.getElementById('favBtn').classList.add('active');
            window.location.href = '/favorites';
            break;
    }
}

       
        async function loadBreedsContent() {
            try {
                const response = await fetch('/breeds/get');
                if (!response.ok) throw new Error('Failed to load breeds');
                
                const breeds = await response.json();
                const container = document.getElementById('breedSearchContainer');
                container.innerHTML = '<iframe src="/breeds" style="width:100%;height:600px;border:none;"></iframe>';
            } catch (error) {
                console.error('Error loading breeds:', error);
            }
        }

       function changeUrl(newPath) {
    setActiveButton(newPath);

    // If the user clicks the breeds button, navigate directly to /breeds
    if (newPath === 'breeds') {
        window.location.href = '/breeds'; // This will trigger the Beego router to load the /breeds page
    } else {
        const newUrl = `/cat/${newPath}`;
        window.history.pushState({ path: newPath }, '', newUrl);
    }
}

        // Event listeners for navigation buttons
        document.getElementById('votingBtn').addEventListener('click', () => changeUrl('voting'));
        document.getElementById('breedBtn').addEventListener('click', () => changeUrl('breeds'));
        document.getElementById('favBtn').addEventListener('click', () => changeUrl('favs'));

        // Handle browser navigation
        window.addEventListener('popstate', () => handleCurrentPath());
    </script>
</body>
</html>
