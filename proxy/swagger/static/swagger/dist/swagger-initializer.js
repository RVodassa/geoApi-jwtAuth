window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  window.ui = SwaggerUIBundle({
    url: "/swagger.yaml",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};

// Проверка загрузки документации
fetch("/swagger.yaml")
  .then(response => {
    if (!response.ok) {
      throw new Error("Ошибка загрузки swagger.yaml");
    }
    return response.json();
  })
  .catch(error => {
    console.error("Ошибка:", error);
  });
