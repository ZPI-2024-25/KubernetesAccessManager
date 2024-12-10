## For English Press Here

# Podręcznik użytkownika

## Wymagania wstępne
1. Zainstalowana i skonfigurowana aplikacja

## 1. Logowanie
Uruchamiamy aplikację webową. Po załadowaniu powinien być widoczny następujący ekran:
![img.png](images/user_guide_1.png)

Po kliknięciu przycisku "Login" zostaniemy przekierowani do strony skonfigurowanego dostawcy tożsamości. Po zalogowaniu zostaniemy przekierowani z powrotem do aplikacji.
![img.png](images/user_guide_2.png)

> *Przykład ekranu logowania na podstawie Keycloak*

## 2. Strona główna
Po zalogowaniu zostaniemy przekierowani do strony głównej aplikacji. Po lewej stronie znajduje się menu, które pozwala na nawigację po aplikacji. Nagłówek zawiera nazwę użytkownika oraz przycisk wylogowania.
![img.png](images/user_guide_3.png)

Kliknięcie przycisku w lewym dolnym rogu pozwala zminimalizować menu. Kliknięcie przycisku ponownie przywraca menu do pełnego rozmiaru.
![img.png](images/user_guide_4.png)

Menu w stanie zminimalizowanym nadal pozwala na nawigację po aplikacji. Po najechaniu kursorem myszy na przyciski menu, wyświetlają się nazwy przycisków. Jeżeli przycisk zawiera podmenu, to po kliknięciu na przycisk zostanie wyświetlone podmenu.
![img.png](images/user_guide_5.png)

## 3. Zarządzanie zasobami na klastrze
Aplikacja pozwala na zarządzanie następującymi zasobami podzielonymi na kategorie:
- Nodes
- Workloads
    - Pods
    - Deployments
    - Daemon Sets
    - Stateful Sets
    - Replica Sets
    - Jobs
    - Cron Jobs
- Config
    - Config Maps
    - Secrets
- Network
    - Services
    - Ingresses
- Storage
    - Persistent Volume Claims
    - Persistent Volumes
    - Storage Classes
- Namespaces
- Access Control
    - Service Accounts
    - Cluster Roles
    - Cluster Role Bindings
- Custom Resources
    - Custom Resource Definitions

### 3.1. Przeglądanie zasobów
Aby wyświetlić listę zasobów, należy wybrać typ zasobu z menu. Po wybraniu typu, sprawdzone zostaną uprawnienia użytkownika do wyświetlenia zasobów. W przypadku braku uprawnień, zasoby nie zostaną wyświetlone.
