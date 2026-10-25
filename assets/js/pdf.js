// HSC Result Checker - PDF Export

function downloadResultPDF() {
    const result = document.getElementById("resultSection");

    if (!result || !result.classList.contains("show")) {
        alert("Please check a result first.");
        return;
    }

    const previousTitle = document.title;

    document.title = "HSC-Result-Checker-Demo-Marksheet";

    window.print();

    document.title = previousTitle;
}

document.addEventListener("DOMContentLoaded", () => {
    const resultContent = document.querySelector(".result-content");

    if (!resultContent) return;

    const pdfButton = document.createElement("button");

    pdfButton.type = "button";
    pdfButton.className = "pdf-btn";
    pdfButton.textContent = "📄 Download PDF";

    pdfButton.addEventListener("click", downloadResultPDF);

    resultContent.appendChild(pdfButton);
});