## For English Press Here

# Podręcznik użytkownika

---

## Wymagania wstępne
1. Zainstalowana i skonfigurowana aplikacja

---

## 1. Logowanie
Uruchamiamy aplikację webową. Po załadowaniu powinien być widoczny następujący ekran:
![Ekran logowania](images/user_guide/user_guide_1.png)

Po kliknięciu przycisku "Login" zostaniemy przekierowani do strony skonfigurowanego dostawcy tożsamości. Po zalogowaniu zostaniemy przekierowani z powrotem do aplikacji.
![Ekran logowania Keycloak](images/user_guide/user_guide_2.png)

> *Przykład ekranu logowania na podstawie Keycloak* 

---

## 2. Strona główna
Po zalogowaniu zostaniemy przekierowani do strony głównej aplikacji. Po lewej stronie znajduje się menu, które pozwala na nawigację po aplikacji. Nagłówek zawiera nazwę użytkownika oraz przycisk wylogowania.
![Strona główna](images/user_guide/user_guide_3.png)

Kliknięcie przycisku w lewym dolnym rogu pozwala zminimalizować menu. Kliknięcie przycisku ponownie przywraca menu do pełnego rozmiaru.
![Zwinięta strona główna](images/user_guide/user_guide_4.png)

Menu w stanie zminimalizowanym nadal pozwala na nawigację po aplikacji. Po najechaniu kursorem myszy na przyciski menu, wyświetlają się nazwy przycisków. Jeżeli przycisk zawiera podmenu, to po kliknięciu na przycisk zostanie wyświetlone podmenu.
![Menu w zwiniętej stronie głównej](images/user_guide/user_guide_5.png)

---

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

---

### 3.1. Przeglądanie zasobów
Aby wyświetlić listę zasobów, należy wybrać typ zasobu z menu. Po wybraniu typu, sprawdzone zostaną uprawnienia użytkownika do wyświetlenia zasobów. W przypadku braku uprawnień, zasoby nie zostaną wyświetlone.
Zasoby wyświetlane są w formie tabeli. Każdy wiersz tabeli reprezentuje jeden zasób. Możliwe jest wybranie liczby zasobów wyświetlanych na stronę. Domyślnie wyświetlanych jest 10 zasobów na stronę.

![Przykładowa lista zasobów](images/user_guide/user_guide_6.png)

W lewym rogu, nad tabelą, znajdują się pola do filtrowania zasobów. Po wpisaniu tekstu w pole search, pozostaną tylko te zasoby, które zawierają wpisany tekst w nazwie. 
Drugie pole pozwala na wybranie namespace, z którego chcemy wyświetlić zasoby. 

Domyślnie wyświetlane są wszystkie zasoby ze wszystkich namespace'ów.

![Zasoby filtrowane po namespace](images/user_guide/user_guide_7.png)
![Zasoby filtrowane po nazwie i namespace](images/user_guide/user_guide_8.png)

W przypadku zasobów nieposiadających namespace, pole namespace nie będzie wyświetlane.

![Przykładowa lista zasobów nie namespace'owanych](images/user_guide/user_guide_9.png)

---

### 3.2. Tworzenie zasobów
Zasób można utworzyć, klikając przycisk "Add" znajdujący się w prawym dolnym rogu. Po kliknięciu przekierowani zostaniemy do edytora tekstu z domyślną konfiguracją zasobu. 
W przypadku braku uprawnień do tworzenia zasobów, przycisk "Add" będzie nieaktywny.

Po zakończeniu edycji należy kliknąć przycisk "Save" w prawym górnym rogu. Zasób zostanie utworzony na klastrze.
Aby porzucić zmiany, należy kliknąć przycisk "Back" w prawym górnym rogu.

![Widok tworzenia zasobów](images/user_guide/user_guide_10.png)

Możliwa jest zmiana formatu definicji zasobu. Domyślnie wyświetlany jest format YAML. Dostępny jest również format JSON.
Pole wyboru formatu znajduje się w lewym górnym rogu. Po wybraniu formatu edytor tekstu zostanie zaktualizowany.

![Definicja zasobu w formacie json](images/user_guide/user_guide_11.png)

Jeżeli zasób jest namespace'owany, to widoczne będzie pole do wpisania namespace'u. 
W przypadku zasobów nieposiadających namespace, pole nie będzie wyświetlane. Domyślnie wybrany jest namespace "default".

![Widok tworzenia zasobu nie namespace'owanego](images/user_guide/user_guide_12.png)

---

### 3.3. Edycja zasobów
Zasób można edytować, klikając w pierwszą ikonę w kolumnie "Actions". Po kliknięciu przekierowani zostaniemy do edytora tekstu z konfiguracją zasobu.
W przypadku braku uprawnień do edycji zasobów, przycisk "Edit" będzie nieaktywny.

Po zakończeniu edycji należy kliknąć przycisk "Save" w prawym górnym rogu. Zasób zostanie zaktualizowany na klastrze.
Aby porzucić zmiany, należy kliknąć przycisk "Back" w prawym górnym rogu.

![Widok edycji zasobu](images/user_guide/user_guide_13.png)

Możliwa jest zmiana formatu definicji zasobu. Domyślnie wyświetlany jest format YAML. Dostępny jest również format JSON.
Pole wyboru formatu znajduje się w lewym górnym rogu. Po wybraniu formatu edytor tekstu zostanie zaktualizowany.

![Widok edycji zasobu w formacie json](images/user_guide/user_guide_14.png)

---

### 3.4. Usuwanie zasobów
Zasób można usunąć, klikając w drugą ikonę w kolumnie "Actions". Po kliknięciu zostanie wyświetlone okno dialogowe z potwierdzeniem usunięcia zasobu.
W przypadku braku uprawnień do usuwania zasobów, przycisk "Delete" będzie nieaktywny.

Po potwierdzeniu zasób zostanie usunięty z klastra.

![Okienko potwierdzające usunięcie zasobu](images/user_guide/user_guide_15.png)

---

### 3.5. Przeglądanie szczegółów zasobu
Szczegóły zasobu można przeglądać, klikając w jakiekolwiek pole wiersza tabeli z wyjątkiem kolumny "Actions". 
Po kliknięciu z prawej strony ekranu wysunie się panel z informacjami o zasobie. Wybrany zasób będzie podświetlony na niebiesko.

Informacja reprezentowana jest na dwa sposoby: tekst lub rozwijane pole z kolejnymi szczegółami.

![Szczegóły zasobu](images/user_guide/user_guide_16.png)

---

## 4. Zarządzanie aplikacjami Helmowymi
Aplikacja pozwala na zarządzanie aplikacjami Helmowymi. W celu zarządzania aplikacjami należy wybrać przycisk "Helm" z menu.

---

### 4.1. Przeglądanie aplikacji Helmowych
Aby wyświetlić listę aplikacji Helmowych, należy wybrać przycisk "Helm" z menu. Po wybraniu przycisku, sprawdzone zostaną uprawnienia użytkownika do wyświetlenia aplikacji helmowych. 
W przypadku braku uprawnień, aplikacje nie zostaną wyświetlone. Możliwe jest wybranie liczby aplikacji wyświetlanych na stronę. Domyślnie wyświetlanych jest 10 aplikacji na stronę.

![Lista aplikacji helmowych](images/user_guide/user_guide_17.png)

W lewym rogu, nad tabelą, znajdują się pola do filtrowania aplikacji. Po wpisaniu tekstu w pole search, pozostaną tylko te aplikacje, które zawierają wpisany tekst w nazwie.
Drugie pole pozwala na wybranie namespace, z którego chcemy wyświetlić aplikacje.

Domyślnie wyświetlane są wszystkie aplikacje ze wszystkich namespace'ów.

---

### 4.2. Przywracanie poprzednich wersji aplikacji Helmowych
Aplikacja pozwala na przywrócenie poprzednich wersji aplikacji Helmowych. Aby przywrócić poprzednią wersję, należy kliknąć przycisk "Restore" w kolumnie "Actions".
W przypadku braku uprawnień do przywracania poprzednich wersji, przycisk będzie nieaktywny.

Po kliknięciu zostanie wyświetlone okno dialogowe z polem do wpisania numeru wersji. Po wpisaniu numeru wersji i potwierdzeniu, zostanie przywrócona poprzednia wersja aplikacji.

![Okienko rollback](images/user_guide/user_guide_18.png)

Jeżeli operacja potrwa dłużej niż 5 sekund, okno zostanie zamknięte, a operacja zostanie wykonana w tle.

---

### 4.3. Odinstalowywanie aplikacji Helmowych
Aplikację Helmową można odinstalować, klikając w drugą ikonę w kolumnie "Actions". Po kliknięciu zostanie wyświetlone okno dialogowe z potwierdzeniem odinstalowania aplikacji.
W przypadku braku uprawnień do odinstalowania aplikacji, przycisk będzie nieaktywny.

Po potwierdzeniu aplikacja zostanie odinstalowana z klastra.

![Okienko potwierdzające odinstalowanie aplikacji](images/user_guide/user_guide_19.png)

---

### 4.4. Przeglądanie szczegółów aplikacji Helmowej
Szczegóły aplikacji Helmowej można przeglądać, klikając w jakiekolwiek pole wiersza tabeli z wyjątkiem kolumny "Actions".
Po kliknięciu z prawej strony ekranu wysunie się panel z informacjami o aplikacji. Wybrana aplikacja będzie podświetlona na niebiesko.

Dla aplikacji Helmowej szczegóły podzielone są na 2 zakładki: "Release" oraz "History".

![Szczegóły aplikacji helmowej](images/user_guide/user_guide_20.png)

---

## 5. Zarządzanie uprawnieniami
Aplikacja pozwala na zarządzanie poziomem dostępu użytkowników poprzez przypisywanie uprawnień rolom uzyskanym od dostawcy tożsamości.

Służy do tego zakładka "Roles" umieszona u dołu menu. Po wybraniu zakładki, sprawdzone zostaną uprawnienia użytkownika do wyświetlania szczegółów oraz edycji ConfigMapy zawierającej mapowanie ról i uprawnień.

---

### 5.1. Przeglądanie ról
Aby wyświetlić listę ról, należy wybrać przycisk "Roles" z menu. Po kliknięciu sprawdzone zostaną uprawnienia użytkownika do wyświetlenia ról.
W przypadku braku uprawnień niemożliwe będzie kliknięcie przycisku "Roles".

Poniżej kolejno widok bez uprawnień oraz widok z uprawnieniami.

![Brak dostępu do ról](images/user_guide/user_guide_21.png)
![Ekran ról](images/user_guide/user_guide_22.png)

---

### 5.2. Szczegóły roli
Szczegóły roli można rozwinąć, klikając w jakąkolwiek część jej pola. Kliknięcie innej roli spowoduje zwinięcie wcześniej rozwiniętej roli.

Każda z ról dzielić się może na 3 sekcje: 
- "Permitted Operations" zawiera operacje, które użytkownik może wykonać.
- "Denied Operations" zawiera operacje, które użytkownik nie może wykonać.
- "Subroles" zawiera podrole, po których dziedziczy dana rola.

Kategorie nie są wyświetlane, jeżeli nie zawierają żadnych operacji.

![Rola zawierająca wszystkie sekcje](images/user_guide/user_guide_23.png)
![Rola zawierająca tylko sekcje "Permitted Operations"](images/user_guide/user_guide_24.png)

Kliknięcie jednej z podroli w sekcji "Subroles" spowoduje rozwinięcie szczegółów tej podroli i przejście do niej.

