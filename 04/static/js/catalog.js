document
    .getElementById("search")
    .addEventListener(
        "input",
        e => {

            const term =
                e.target.value
                    .toLowerCase();

            document
                .querySelectorAll(".api-card")
                .forEach(card => {

                    const text =
                        card.dataset.search
                            .toLowerCase();

                    card.style.display =
                        text.includes(term)
                        ? "block"
                        : "none";
                });
        }
    );