const hangmanImage = document.getElementById("hangman-image");

function updateImage(attempts) {
  hangmanImage.src = `/static/img/hangman/hang${attempts}.png`;
}

// Exemple d'appel AJAX simulé
fetch("/get-attempts") // Route à définir dans Go pour récupérer les tentatives restantes
  .then((response) => response.json())
  .then((data) => {
    updateImage(data.attempts);
  });
