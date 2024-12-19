document.addEventListener("DOMContentLoaded", function () {
    // Get references to placeholders and sections
    const catImagesContainer = document.getElementById("cat-images");
    const votingSection = document.getElementById("voting-section");
    const breedsSection = document.getElementById("breeds-section");
    const favsSection = document.getElementById("favs-section");

    const votingImage = document.getElementById("voting-image");
    const votingHeart = document.getElementById("voting-heart");
    const votingLike = document.getElementById("voting-like");
    const votingDislike = document.getElementById("voting-dislike");

    const breedsSearch = document.getElementById("breeds-search");
    const breedsSuggestions = document.getElementById("breeds-suggestions");
    const breedImage = document.getElementById("breed-image");
    const breedName = document.getElementById("breed-name");
    const breedOrigin = document.getElementById("breed-origin");
    const breedDescription = document.getElementById("breed-description");
    const breedWiki = document.getElementById("breed-wiki");
    const sliderDots = document.getElementById("slider-dots");

    const favoritesContainer = document.getElementById("favorites-container");

    let favorites = [];
    let activeTab = "voting"; // Keep track of the active tab

    // Function to set active tab styles and display appropriate sections
    function setActiveTab(tabId) {
        const tabs = ["tab-voting", "tab-breeds", "tab-favs"];
        tabs.forEach((tab) => {
            const tabElement = document.getElementById(tab);
            const sectionId = tab.replace("tab-", "") + "-section";
            const sectionElement = document.getElementById(sectionId);

            // Highlight the active tab and show its section
            if (tab === tabId) {
                tabElement.classList.add("text-red-500", "border-b-2", "border-red-500");
                tabElement.classList.remove("text-gray-500");
                sectionElement.style.display = "block";
            } else {
                // Reset inactive tabs and hide their sections
                tabElement.classList.remove("text-red-500", "border-b-2", "border-red-500");
                tabElement.classList.add("text-gray-500");
                sectionElement.style.display = "none";
            }
        });
    }

    // Function to fetch and display one random cat image in the Voting Tab
    async function fetchSingleCatImage() {
        try {
            const response = await fetch("/api/cats");
            const images = await response.json();
            if (images.length > 0) {
                renderSingleCatImage(images[0]); // Pass the first image to render function
            }
        } catch (error) {
            console.error("Error fetching cat image:", error);
        }
    }

    function renderSingleCatImage(image) {
        const placeholderImageUrl = "/static/img/placeholder2.gif"; // Path to placeholder image
    
        // Show placeholder image immediately
        votingImage.src = placeholderImageUrl;
        votingImage.alt = "Loading..."; // Update alt text for accessibility
    
        // Add functionality to the heart button
        votingHeart.onclick = async () => {
            handleFavorite(image); // Add to favorites
            votingImage.src = placeholderImageUrl; // Set placeholder
            await fetchSingleCatImage(); // Fetch and display the next image
        };
    
        // Add functionality to the like button
        votingLike.onclick = async () => {
            votingImage.src = placeholderImageUrl; // Set placeholder
            await fetchSingleCatImage(); // Fetch and display the next image
        };
    
        // Add functionality to the dislike button
        votingDislike.onclick = async () => {
            votingImage.src = placeholderImageUrl; // Set placeholder
            await fetchSingleCatImage(); // Fetch and display the next image
        };
    
        // After fetchSingleCatImage updates with the actual image
        setTimeout(() => {
            votingImage.src = image.url; // Set the actual image
            votingImage.alt = "Cat"; // Update alt text
        }, 100); // Optionally add a small delay for smooth transitions
    }

    // Function to fetch and render breeds data with search, dropdown, and slider functionality
    async function fetchBreedsTab() {
        try {
            // Clear previous content
            const response = await fetch("/api/breeds");
            const breeds = await response.json();

            // Ensure breeds are fetched
            if (breeds.length === 0) {
                console.error("No breeds available.");
                return;
            }

            console.log("Breeds tab active");

            // Automatically load details for the first suggestion
            const firstBreed = breeds[0];
            updateBreedDetails(firstBreed);
            console.log("First selected breed: ", firstBreed.id, firstBreed.name);
            breedImage.alt = "Please select a breed";
    
            // Add event listener for the search bar
            breedsSearch.addEventListener("input", (event) => {
                const query = event.target.value.toLowerCase();
                const filteredBreeds = breeds.filter((breed) =>
                    breed.name.toLowerCase().includes(query)
                );
                console.log(filteredBreeds);
                renderSuggestions(filteredBreeds);
            });
    
            // Render the suggestions dropdown
            function renderSuggestions(filteredBreeds) {
                breedsSuggestions.innerHTML = ""; // Clear previous suggestions
                if (filteredBreeds.length > 0) {
                    breedsSuggestions.style.display = "block";
                    filteredBreeds.forEach((breed) => {
                        const suggestion = document.createElement("div");
                        suggestion.className = "suggestions";
                        suggestion.textContent = breed.name;
    
                        // Handle breed selection
                        suggestion.addEventListener("click", () => {
                            breedsSearch.value = breed.name;
                            breedsSuggestions.style.display = "none";
                            console.log("Selected breed from dropdown: ", breed.name);
                            console.log("Selected breed id from dropdown: ", breed.id);
                            updateBreedDetails(breed);
                        });
    
                        breedsSuggestions.appendChild(suggestion);
                    });
                } else {
                    breedsSuggestions.style.display = "none";
                }
            }
    
            // Update breed details and images
            async function updateBreedDetails(breed) {
                console.log("Updating Breed Details:", breed.id, breed.name);
            
                // Update text content
                breedName.textContent = breed.name || "Unknown Breed";
                breedOrigin.textContent = breed.origin ? `(${breed.origin})` : "Origin not available";
                breedDescription.textContent = breed.description || "Description not available.";
                breedWiki.href = breed.wikipedia_url || "#";
            
                // Fetch images for the selected breed
                const imagesResponse = await fetch(`/api/cats?breed_id=${breed.id}`);
                console.log("Breed id is: ", breed.id);
                const images = await imagesResponse.json();

                //const limitedImages = images.slice(0, 5);
            
                if (images.length > 0) {
                    // Set up a slider
                    setupSlider(images, breed.name);
                } else {
                    breedImage.src = "/static/images/placeholder.webp"; // Fallback placeholder
                    breedImage.alt = "No image available";
                }
            }
        } catch (error) {
            console.error("Error fetching breeds:", error);
        }
    }

    function setupSlider(images, breedName) {
        let currentIndex = 0;
    
        // Function to update the displayed image
        function showImage(index) {
            console.log("Show image function called");
            console.log("Displaying image:", images[index].url);
            breedImage.src = images[index].url;
            breedImage.alt = `${breedName} Image ${index + 1}`;
            updateDots(index);
        }
    
        // Function to update the dots for the slider
        function updateDots(index) {
            console.log("Dots updated");
            sliderDots.innerHTML = ""; // Clear existing dots
            images.forEach((_, i) => {
                const dot = document.createElement("div");
                dot.className = "dot";
                dot.style.backgroundColor = i === index ? "#000000" : "#808080"; 
                dot.addEventListener("click", () => {
                    currentIndex = i;
                    showImage(currentIndex);
                });
                sliderDots.appendChild(dot);
            });
        }
    
        // Initialize the slider with the first image
        showImage(currentIndex);
    }

    // Function to render favorite images in a grid layout
    function renderCatImages(images) {
        favoritesContainer.innerHTML = ""; // Clear existing content
        images.forEach((image) => {
            const imgCard = document.createElement("div");
            imgCard.className = "p-2 border rounded-lg shadow-sm";

            imgCard.innerHTML = `
                <img src="${image.url}" alt="Cat" class="w-full h-48 object-cover rounded-lg mb-2">
                <div class="flex justify-between items-center">
                    <button class="remove-fav-btn text-red-500 font-semibold">Remove</button>
                </div>
            `;

            // Remove from favorites
            imgCard.querySelector(".remove-fav-btn").addEventListener("click", () => {
                favorites = favorites.filter((fav) => fav.id !== image.id);
                renderCatImages(favorites);
            });

            favoritesContainer.appendChild(imgCard);
        });
    }

    // Handle adding an image to favorites
    async function handleFavorite(image) {
        if (!favorites.find((fav) => fav.id === image.id)) {
            favorites.push(image);
            alert("Added to favorites!");
        } else {
            alert("Already in favorites!");
        }
    }

    // Tab Event Listeners
    document.getElementById("tab-voting").addEventListener("click", () => {
        setActiveTab("tab-voting");
        fetchSingleCatImage();
    });

    document.getElementById("tab-breeds").addEventListener("click", () => {
        setActiveTab("tab-breeds");
        fetchBreedsTab();
    });

    document.getElementById("tab-favs").addEventListener("click", () => {
        setActiveTab("tab-favs");
        renderCatImages(favorites);
    });

    // Initialize with the Voting Tab
    setActiveTab("tab-voting");
    fetchSingleCatImage();
});

