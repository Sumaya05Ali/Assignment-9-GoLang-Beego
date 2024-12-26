# Assignment-9-GoLang-Beego
# CAT API Project

This project is a web application built with the Beego framework in Go, designed to provide a fun and interactive experience for users to explore cat images, learn about different breeds, vote on images, and save their favorite cat pictures.

## Features

- Voting on Images: Users can like or dislike cat images.

- Cat Breeds: Browse through various cat breeds and view related information.

- Favorite Images: Save your favorite cat images and view them in a dedicated section.

- Dynamic Navigation: Seamlessly navigate between different sections of the app.


## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [API Integration](#api-integration)
- [Frontend Features](#frontend-features)

## Installation

Prerequisites:

1. Install Go.

2. Install Beego and its CLI tool:

   ```bash
    go get -u github.com/beego/beego/v2
    go install github.com/beego/bee/v2
    ```
3. Install dependencies for the project:
   ```bash
      go mod tidy
    ```

Steps:
 
1. Clone the repository:
   ```bash
   git clone https://github.com/Sumaya05Ali/Assignment-9-GoLang-Beego.git
   ```
2. Navigate to the project directory:
   ```bash
    cd  Assignment-9-GoLang-Beego-main
   ```
3. Run the application:
   ```bash
        bee run
     ```
   The application will be available at http://localhost:8080 by default.

4. Testing:
   ```bash
    go test ./... -v
    go test ./... -v -cover
    go test ./... -v -coverprofile=coverage.out
    go tool cover -html=coverage.out
     ```

## Usage

- Navigate to http://localhost:8080/cat to start viewing cat images.

- Use the navigation bar to:

- Vote on images (like/dislike).

- Browse through cat breeds.

- Save your favorite images and view them in the favorites section.

## API Integration

This application integrates with The Cat API for fetching cat images and breed information. The API_KEY is used to authenticate requests.

Endpoints Used

1. Fetch Cat Images: https://api.thecatapi.com/v1/images/search?limit=1

2. Get Breed Info: Custom API endpoint /breeds/get served by the Beego backend.

## Frontend Features

- Responsive Design: The UI adapts to various screen sizes.

- Interactive Buttons: Buttons for liking, disliking, and saving favorites.

- Dynamic Navigation: Client-side navigation implemented using JavaScript.

- Notifications: Feedback messages for user actions like saving favorites.

## License

This project is licensed under the MIT License. See the LICENSE file for details.  
