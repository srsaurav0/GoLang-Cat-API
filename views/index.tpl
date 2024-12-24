<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Voting App</title>
    <!-- <link href="static/css/output.css"  rel="stylesheet"> -->
    <link href="static/css/global.css"  rel="stylesheet">
    <!-- <link href="https://cdn.jsdelivr.net/npm/tailwindcss@3.3.2/dist/tailwind.min.css" rel="stylesheet"> -->
</head>
<body class="bg-gray-100">

    <!-- Main Container -->
    <div class="max-w-4xl mx-auto my-10 p-4 bg-white shadow-lg rounded-lg">
        <!-- Tabs -->
        <div class="flex space-x-6 border-b pb-2 mb-4">
            <button id="tab-voting" class="font-semibold text-red-500 border-b-2 border-red-500">Voting</button>
            <button id="tab-breeds" class="font-semibold text-gray-500">Breeds</button>
            <button id="tab-favs" class="font-semibold text-gray-500">Favs</button>
        </div>

        <!-- Cat Images Section -->
        <div id="cat-images">
            <!-- Voting Section -->
            <div id="voting-section">
                <div class="p-2 border rounded-lg shadow-sm">
                    <!-- Image -->
                    <img id="voting-image" src="static/img/placeholder2.gif" alt="Cat" class="p_img">
                    
                    <!-- Buttons -->
                    <div class="flex justify-between space-x-6">
                        <button id="voting-heart" style="font-size: 60px; ">❤️</button>
                        <button id="voting-like" style="font-size: 60px;">👍</button>
                        <button id="voting-dislike" style="font-size: 60px;">👎</button>
                    </div>
                </div>
            </div>

            <!-- Breeds Section -->
            <div id="breeds-section" style="display: none;">
                <div class="relative w-full mb-4">
                    <input id="breeds-search" type="text" placeholder="Search breeds..." class="p-2 border rounded-lg w-full">
                    <div id="breeds-suggestions" class="dropdown" style="display: none;"></div>
                </div>
                <div class="w-full text-center">
                    <div class="relative w-full mx-auto mb-4">
                        <img id="breed-image" src="static/img/placeholder2.gif" alt="Cat" class="p_img">
                        <div id="slider-dots" class="flex justify-center items-center"></div>
                        <div class="flex justify-center mt-4">
                            <button id="left-button" class="text-white px-4 py-2 rounded hover:bg-gray-200 text-6xl" >⬅️</button>
                            <button id="right-button" class="text-white px-4 py-2 rounded hover:bg-gray-200 text-6xl">➡️</button>
                        </div>
                    </div>
                    <h2 id="breed-name" class="text-2xl font-bold mt-4"></h2>    
                    <div class="flex justify-center mt-4">
                        <p id="breed-origin" class="text-gray-500 font-bold"></p>
                        <p id="breed-id" class="text-gray-500 italic ml-2"></p>
                    </div>
                    <p id="breed-description" class="mt-2"></p>
                    <a id="breed-wiki" href="#" class="text-blue-500 underline">WIKIPEDIA</a>
                </div>                
            </div>

            <!-- Favorites Section -->
            <div id="favs-section" style="display: none;">
                <div id="favorites-container" class="grid grid-cols-2 gap-4"></div>
            </div>
        </div>
    </div>

    <script src="/static/js/app.js"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</body>
</html>
