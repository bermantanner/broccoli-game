cards = document.querySelectorAll('.game-card');

cards.forEach(card => {
    card.addEventListener('click', function() {
        if (card.classList.contains('selected')) {
            card.classList.remove('selected');
        }
        else {
            cards.forEach(c => {
                c.classList.remove('selected');
            });
            card.classList.add('selected');
        }
    });
});