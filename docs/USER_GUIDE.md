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
Zasoby wyświetlane są w formie tabeli. Każdy wiersz tabeli reprezentuje jeden zasób. Możliwe jest wybranie liczby zasobów wyświetlanych na stronę. Domyślnie wyświetlanych jest 10 zasobów na stronę.

![img.png](images/user_guide_6.png)

W lewym rogu, nad tabelą, znajdują się pola do filtrowania zasobów. Po wpisaniu tekstu w pole search, pozostaną tylko te zasoby, które zawierają wpisany tekst w nazwie. 
Drugie pole pozwala na wybranie namespace, z którego chcemy wyświetlić zasoby. 

Domyślnie wyświetlane są wszystkie zasoby ze wszystkich namespace'ów.

![img.png](images/user_guide_7.png)
![img.png](images/user_guide_8.png)

### 3.2. Tworzenie zasobów
Zasób można utworzyć klikając przycisk "Add" znajdujący się w prawym dolnym rogu. Po kliknięciu przekierowani zostaniemy do edytora tekstu z domyślną konfiguracją zasobu. 
Po zakończeniu edycji należy kliknąć przycisk "Save" w prawym górnym rogu. Zasób zostanie utworzony na klastrze.

Aby porzucić zmiany, należy kliknąć przycisk "Back" w prawym górnym rogu.

