<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Breeds</title>
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
            margin: 0 auto;
            padding: 20px;
        }

        .search-container {
    margin-bottom: 20px;
    position: relative;
    padding: 0 10px;
}

        .search-input {
    width: 92%;
    padding: 12px 20px;
    border: 2px solid #eee;
    border-radius: 25px;
    font-size: 16px;
    transition: all 0.3s ease;
    background-color: #f8f8f8;
}

       .search-input:focus {
    outline: none;
    border-color: #FF4D4D;
    background-color: white;
    box-shadow: 0 2px 8px rgba(255, 77, 77, 0.1);
}

        .breeds-list {
    position: absolute;
    top: calc(100% + 5px);
    left: 10px;
    right: 10px;
    background: white;
    border: 1px solid #eee;
    border-radius: 15px;
    max-height: 300px;
    overflow-y: auto;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    z-index: 1000;
    margin-top: 5px;
    display: none;
}

       .breed-option {
    padding: 12px 20px;
    cursor: pointer;
    transition: 0.2s;
    display: flex;
    justify-content: space-between;
    align-items: center;
}
       .breed-option:hover {
    background-color: #f8f8f8;
}
  .breed-option .breed-name {
    font-size: 16px;
    margin: 0;
    color: #333;
}

.breed-option .breed-id {
    font-size: 14px;
    color: #999;
    font-family: monospace;
}
        .breed-content {
            padding: 20px 0;
        }

        .breed-image {
            width: 100%;
            height: 400px;
            object-fit: cover;
            border-radius: 15px;
            margin-bottom: 20px;
        }

        .breed-info {
            text-align: left;
        }

        .breed-name {
            font-size: 24px;
            margin-bottom: 10px;
            color: #333;
        }

        .breed-origin {
            color: #666;
            margin-bottom: 15px;
            font-style: italic;
        }

        .breed-description {
            color: #444;
            line-height: 1.6;
            margin-bottom: 20px;
        }

        .wiki-link {
            display: inline-block;
            padding: 10px 20px;
            background-color: #FF4D4D;
            color: white;
            text-decoration: none;
            border-radius: 8px;
            transition: 0.3s;
        }

        .wiki-link:hover {
            background-color: #ff3333;
        }

        .image-navigation {
            display: flex;
            justify-content: center;
            gap: 8px;
            margin: 15px 0;
        }

        .dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background-color: #ddd;
            cursor: pointer;
            transition: 0.3s;
        }

        .dot.active {
            background-color: #FF4D4D;
        }

        .nav-buttons {
            display: flex;
            gap: 15px;
            margin-bottom: 20px;
        }

        .nav-button {
            padding: 10px 20px;
            border: none;
            background: none;
            cursor: pointer;
            color: #666;
            font-size: 16px;
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .nav-button.active {
            color: #FF4D4D;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="nav-buttons">
            <button class="nav-button" onclick="window.location.href='/cat/voting'">
                <i class="fas fa-arrow-up"></i>
                Voting
            </button>
            <button class="nav-button active">
                <i class="fas fa-search"></i>
                Breeds
            </button>
            <button class="nav-button" onclick="window.location.href='/cat/favs'">
                <i class="fas fa-heart"></i>
                Favs
            </button>
        </div>

        <div class="search-container">
            <input type="text" class="search-input" placeholder="Search for breeds...">
            <div class="breeds-list"></div>
        </div>

        <div class="breed-content"></div>
    </div>

    <script>
        let breeds = [];
        let currentImageIndex = 0;
        let currentBreed = null;
        let images = [];
        const searchInput = document.querySelector('.search-input');
        const breedsList = document.querySelector('.breeds-list');
        const breedContent = document.querySelector('.breed-content');

        async function fetchBreeds() {
            try {
                const response = await fetch('/breeds/get');
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                breeds = await response.json();
                showBreedsList();
            } catch (error) {
                console.error('Error fetching breeds:', error);
            }
        }

       function showBreedsList() {
    const searchTerm = searchInput.value.toLowerCase();
    const filteredBreeds = breeds.filter(breed =>
        breed.name.toLowerCase().includes(searchTerm)
    );

    if (filteredBreeds.length === 0) {
        breedsList.style.display = 'none';
        return;
    }

    breedsList.innerHTML = filteredBreeds
        .map(breed => `
            <div class="breed-option" data-id="${breed.id}">
                <span class="breed-name">${breed.name}</span>
                <span class="breed-id">${breed.id}</span>
            </div>
        `)
        .join('');
    breedsList.style.display = 'block';
}
       async function displayBreed(breed) {
    currentBreed = breed;
    currentImageIndex = 0;
    images = [];

    // First, add the default image from the breed data
    if (breed.image && breed.image.url) {
        images.push(breed.image.url);
    }

    // Then try to fetch additional images
    try {
        const imageResponse = await fetch(`/images/${breed.id}`);
        
        if (imageResponse.ok) {
            const additionalImages = await imageResponse.json();
            // Add additional image URLs to the images array
            additionalImages.forEach(img => {
                if (img.url && !images.includes(img.url)) {
                    images.push(img.url);
                }
            });
        }
    } catch (error) {
        console.warn('Could not fetch additional images:', error);
    }

    // If no images were found, use the default image or a placeholder
    if (images.length === 0) {
        images = ['/api/placeholder/600/400'];
    }

    updateBreedDisplay();
}

        function updateBreedDisplay() {
            if (!currentBreed || images.length === 0) return;

            const dotsHtml = images.length > 1 
                ? images.map((_, i) => 
                    `<div class="dot ${i === currentImageIndex ? 'active' : ''}" data-index="${i}"></div>`
                ).join('')
                : '';

            breedContent.innerHTML = `
                <img src="${images[currentImageIndex]}" alt="${currentBreed.name}" class="breed-image">
                ${images.length > 1 ? `<div class="image-navigation">${dotsHtml}</div>` : ''}
                <div class="breed-info">
                    <h2 class="breed-name">${currentBreed.name}</h2>
                    <p class="breed-origin">${currentBreed.origin}</p>
                    <p class="breed-description">${currentBreed.description}</p>
                    <a href="https://wikipedia.org/wiki/${currentBreed.name.replace(/\s+/g, '_')}" 
                       class="wiki-link" target="_blank">WIKIPEDIA</a>
                </div>
            `;

            // Add click handlers for dots
            document.querySelectorAll('.dot').forEach(dot => {
                dot.addEventListener('click', () => {
                    currentImageIndex = parseInt(dot.dataset.index);
                    updateBreedDisplay();
                });
            });
        }

        // Auto-rotate images only if there are multiple images
        setInterval(() => {
            if (currentBreed && images.length > 1) {
                currentImageIndex = (currentImageIndex + 1) % images.length;
                updateBreedDisplay();
            }
        }, 3000);

        // Event Listeners
        searchInput.addEventListener('focus', () => {
            if (breeds.length === 0) {
                fetchBreeds();
            } else {
                showBreedsList();
            }
        });

        searchInput.addEventListener('input', showBreedsList);

        breedsList.addEventListener('click', e => {
            const option = e.target.closest('.breed-option');
            if (option) {
                const breed = breeds.find(b => b.id === option.dataset.id);
                if (breed) {
                    displayBreed(breed);
                    searchInput.value = breed.name;
                    breedsList.style.display = 'none';
                }
            }
        });

        document.addEventListener('click', e => {
            if (!e.target.closest('.search-container')) {
                breedsList.style.display = 'none';
            }
        });

        // Initial load
        fetchBreeds();
    </script>
</body>
</html>