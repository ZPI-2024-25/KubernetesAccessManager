## Opis autoryzacji w aplikacji
### Token JWT
Aplikacja zakłada połączenie z dostarczycielem tożsamości Keycloak lub innym spełniającym wymagania protokołu OpenIdConnect i dostarczającego role. Autoryzacja oparta jest na ekstrakcji ról użytkownika z tokena JWT wysyłanego z każdym zapytaniem do serwera backend. Aplikacja odczytuje role zdefiniowane w dwóch miejscach tokena:
- resource_access → `VITE_KEYCLOAK_CLIENTNAME` → roles: Lista ról przypisanych w kontekście konkretnego klienta.
- realm_access → roles: Lista ról przypisanych globalnie w ramach danego realm.

Wartość `VITE_KEYCLOAK_CLIENTNAME` odpowiada nazwie klienta skonfigurowanej jako zmienna środowiskowa w Helm Chart. Zmienna ta wskazuje na sekcję tokena odpowiadającą aplikacji.
```json
{
  "exp": 1733164880,
  "iat": 1733164580,
  "auth_time": 1733164580,
  "jti": "ed620242-9507-4ab7-bdef-d7e0c5c132e4",
  "iss": "http://localhost:4000/realms/ZPI-realm",
  "aud": "account",
  "sub": "2dffab05-b413-40bd-82e1-bbdd9639d5b4",
  "typ": "Bearer",
  "azp": "ZPI-client",
  "sid": "e355c61f-5394-4110-98ac-8759c6012596",
  "acr": "1",
  "allowed-origins": [
    "*"
  ],
  "realm_access": {
    "roles": [
      "default-roles-zpi-realm",
      "realm-zpi-role",
    ]
  },
  "resource_access": {
    "ZPI-client": {
      "roles": [
        "zpi-role"
      ]
    },
    "account": {
      "roles": [
        "manage-account",
        "manage-account-links",
        "view-profile"
      ]
    }
  },
  "scope": "openid email profile",
  "email_verified": false,
  "name": "a a",
  "preferred_username": "zpi-user",
  "given_name": "a",
  "family_name": "a",
  "email": "a@a.a"
}
```
W tym wypadku dla `VITE_KEYCLOAK_CLIENTNAME` równego `"ZPI-client"` aplikacja by autoryzowała żądanie na podstawie ról `"zpi-role"`, `"default-roles-zpi-realm"`, `"realm-zpi-role"`. Jeśli przynajmniej jedna z wczytanych ról daje użytkownikowi dostęp do określonego zapytania, to użytkownik jest pomyślnie autoryzowany. W przeciwnym wypadku zwracany jest błąd 403 - jeżeli żadna z rola nie nadaje mu odpowiednich uprawnień.

### Możliwe uprawnienia w aplikacji
W aplikacji wyróżniamy akcje create, read, update, delete, list. Każde z uprawnień akcji musi być nadane osobno, w szczególności uprawnienie do listowania zasobów nie daje automatycznie możliwości przeglądania ich szczegółów. Uprawnienia można definiować dla konkretnej przestrzeni nazw i typu zasobu np. Pod, ConfigMap.

### Mapa ról na uprawnienia
W celu korzystania z aplikacji należy przed instalacją zdefiniować zasób typu ConfigMap z kluczami role-map i subrole-map(opcjonalnie). Mapę ról można zmieniać w trakcie działania aplikacji. Nazwę i namespace ConfigMap’y należy zdefiniować w zmiennych środowiskowych `ROLEMAP_NAMESPACE` i `ROLEMAP_NAME` . Domyślne wartości to namespace "default" i "role-map". Przykładowa definicja:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: role-map 
  namespace: default
data:
  role-map: |
    superadmin:
      permit:
        - operations: ["*"]
    admin:
      deny: 
        - namespace: "top-restricted"
        - namespace: "role-map-namespace"
          resource: "ConfigMap"
          operations: ["delete", "create", "update"]
      permit:
        - operations: ["*"]
    manager:
      deny:
        - ["delete", "create", "update"]
      subroles:
        - "admin1"
        - "admin2"
        - "permissionsViewer"
    team1Admin: 
      subroles:
        - "kubeConfigViewer"
        - "team1Admin"
        - "permissionsViewer"
    team2Admin:
      subroles:
        - "kubeConfigViewer"
        - "team2Admin"
        - "permissionsViewer"
  subrole-map: |
    team1Admin:
      permit:
        - namespace: "team1"
    team2Admin:
      permit:
        - namespace: "team2"
    kubeConfigViewer:
      permit:
        - namespace: "kube-system"
          operations: ["read", "list"]
      deny:
        - resource: "secretResource"
    permissionsViewer:
      permit:
        - namespace: "role-map-namespace"
          resource: "ConfigMap"
          operations: ["read", "list"]
```
Pola definicji roli i podroli są takie same. Są to:

- **nazwa roli/podroli** - podana przed uprawnieniami, jako klucz. W przypadku roli powinna być dokładnie taka sama jak rola odczytana z listy ról w tokenie JWT.
- **permit** - lista operacji na które zezwala dana rola
- **deny** - lista operacji które rola zabrania
- **subroles** - lista nazw podról, z których dziedziczone są uprawnienia.

Oprócz samej nazwy roli/podroli należy zdefiniować przynajmniej jeden atrybut roli(permit/deny/subroles).

### Definiowanie operacji w permit, deny
Definiując operacje, można uwzględnić 3 atrybuty - namespace, resource(typ zasobu np. “Pod”) oraz akcję (CRUDL). Każde z nich można pominąć co będzie jednoznaczne z ograniczeniem/nadaniem uprawnień dla wszystkich możliwych namespace'ów, resource'ów lub akcji. Taki sam efekt można uzyskać wpisując “*” dla namespace i resource lub [“*”] dla operations. Przykład
```yaml
    admin:
      deny: 
        - namespace: "top-restricted"
        - namespace: "role-map-namespace"
          resource: "ConfigMap"
          operations: ["delete", "create", "update"]
      permit:
        - operations: ["*"]
```
Według powyższej definicji roli, admin ma prawo do wszystkich akcji, na wszystkich zasobach w dowolnym namespace, za wyjątkiem namespace “top-restricted” oraz usuwania, edytowania lub tworzenia ConfigMap w namespace `“role-map-namespace”`. Wartość `operations: [“*”]` jest wymagana, gdyż potrzebny jest przynajmniej jeden atrybut z namespace, resource, operations.

### Używanie podról
Używając podról można zdefiniować konfiguracje uprawnień często powtarzających się pomiędzy poszczególnymi rolami. Ważne jest rozróżnienie pomiędzy rolą a podrolą - nazwa roli pochodzi z zewnętrznego dostawcy tożsamości i musi być dokładnie taka sama jak w tokenie JWT aby gwarantować użytkownikowi jakiekolwiek uprawnienia. Podrola służy wyłącznie w przekazywaniu uprawnień roli. Można zdefiniować rolę i podrolę o tej samej nazwie. Aby rola otrzymała uprawnienia z podroli, należy dodać nazwę podroli do listy “subroles” w roli. Nie można używać ról jako podról. Podrole też mogą posiadać podrole.

Ograniczenia (deny) z podroli nie wpływają na uprawnienia (permit) zdefiniowane w nadroli oraz na uprawnienia wynikające z pozostałych podról nadroli. Natomiast ograniczenia (deny) z nadroli wpływają na uprawnienia (permit) ze wszystkich podroli.

Przykład:
```yaml
  role-map: |
    manager:
      deny:
        - ["delete", "create", "update"]
      subroles:
        - "team1admin"
        - "team2admin"
    team1admin: 
      subroles:
        - "team1admin"
    team2Admin:
      subroles:
        - "team2admin"
  subrole-map: |
    team1admin:
      permit:
        - namespace: "team1"
      subroles:
        - "permissionViewer"
    team2admin:
      permit:
        - namespace: "team2"
      subroles:
        - "permissionViewer"
    permissionsViewer:
      permit:
        - namespace: "role-map-namespace"
          resource: "ConfigMap"
          operations: ["read", "list"]
```
W przedstawionym przykładzie `team1Admin` jest jednocześnie zdefiniowany jako rola i podrola. Nie powoduje to problemu, gdyż role i podrole są zawsze rozpatrywane jako osobne byty. Podrola `team1admin` dziedziczy uprawnienia z `permissionViewer`, umożliwiając dowolne operacje w namespace `team1` oraz odczyt i listowanie ConfigMap w namespace `"role-map-namespace"`. Analogicznie działa `team2admin`, operując w namespace `"team2"`. Podrole te są podrolami dla ról `team1admin`, `team2admin` oraz dla `manager`.

Rola `manager` niezależnie od uprawnień podról (`team1admin`, `team2admin`) nie może wykonywać operacji delete, create ani update, co wynika z nadrzędnych ograniczeń (deny). Ograniczenia te mają wyższy priorytet niż uprawnienia z podról. Podrole mogą dziedziczyć uprawnienia od innych podról, co widać na przykładzie `team1Admin` i `permissionViewer`.

Ostatecznie, rola `team1Admin` gwarantuje uprawnienia do wszystkich zasobów i akcji w namespace `"team1"` oraz podgląd i listowanie ConfigMap w namespace `"role-map-namespace"`, analogicznie `team2admin` dla namespace `"team2"`. Rola `manager` może listować i czytać szczegóły zasobów w namespace’ach `"team1"`, `"team2"` oraz ConfigMap’y w namespace `"role-map-namespace"`