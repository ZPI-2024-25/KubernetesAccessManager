<h1 style="display: flex; align-items: center; gap: 30px;">
  Kubernetes Access Manager (KAM)
  <img src="./docs/images/KAM_logo.png" alt="KAM Logo" style="height: 125px;">
</h1>

---
## For English Press [Here](#english) / Для русского нажмите [здесь](#russian)

---

Aplikacja **Kubernetes Access Manager (KAM)** pomaga rozwiązać problem braku elastyczności w zarządzaniu dostępem do klastrów Kubernetes. Dzięki precyzyjnemu przypisywaniu uprawnień na podstawie przestrzeni nazw (namespace) i typów zasobów (resource), KAM upraszcza zarządzanie klastrami i zwiększa bezpieczeństwo.

Zintegrowana z systemami zarządzania tożsamością, takimi jak Keycloak, aplikacja wykorzystuje standard OpenID Connect (OIDC), zapewniając łatwą autoryzację i uwierzytelnianie.

---

## Najważniejsze funkcje

- **Zarządzanie zasobami Kubernetes**: Tworzenie, edycja, usuwanie.
- **Przegląd zasobów**:
    - Wyświetlanie listy zasobów.
    - Szczegółowe informacje o zasobach i aplikacjach helmowych.
- **Rollback aplikacji helmowych**: Przywracanie wcześniejszych wersji.
- **Zarządzanie dostępem**:
    - 5 poziomów uprawnień:
        - Widok listy.
        - Widok szczegółowy.
        - Tworzenie.
        - Edycja.
        - Usuwanie.
- **Integracja z Keycloak**: Łatwa autoryzacja i uwierzytelnianie.
- **Intuicyjny interfejs webowy**: Dostosowany do użytkowników o różnym poziomie zaawansowania.
- **Prosta instalacja**: Możliwość wdrożenia w Kubernetes za pomocą Helm Chart.

---

## Diagram wysokopoziomowy

### Architektura aplikacji
![Diagram architektury Helm Chart](./path/to/helm_chart.png)

### Przykład działania: CRUD Kubernetes
![Diagram sekwencji CRUD Kubernetes](./path/to/diagram_crud_k8s.png)

---

## Jak zacząć?

1. **Zainstaluj aplikację za pomocą Helm Chart**:
   ```bash
   helm repo add kam-repo https://repo.example.com/charts
   helm install kam kam-repo/kubernetes-access-manager
   ```

2. **Uruchom backend i frontend** zgodnie z instrukcją w [quickstart.md](./path/to/quickstart.md).

3. **Zaloguj się** w aplikacji i zacznij zarządzać klastrami Kubernetes.

---

## Dokumentacja

- [Quickstart](./path/to/quickstart.md)
- [Instrukcja wdrożenia](./path/to/deployment.md)
- [Dokumentacja API](./path/to/api.md)
- [Diagramy](./path/to/diagrams/)

---

## Autorzy

- **Vera Goriukhina** – Frontend
- **Marek Fiuk** – Napiszę się tu jeszcze dokładnie
- **Dawid Walkiewicz** – Pan maruda junior aka łamliwy kark 
- **Samuel Żołądz** – Pan maruda senior, obrońca dobrze zrobionych JSON'ów

Jeśli masz pytania lub sugestie, skontaktuj się z nami poprzez [issues](https://github.com/ZPI-2024-25/KubernetesAccessManager/issues).

---

## Licencja

Projekt jest dostępny na licencji Beerware (się ustali). Szczegóły w pliku [LICENSE](./LICENSE).
