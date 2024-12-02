## For english press [here](#authorization-in-the-application)
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
W tym wypadku dla `VITE_KEYCLOAK_CLIENTNAME` równego `"ZPI-client"` aplikacja autoryzowałaby żądanie na podstawie ról `"zpi-role"`, `"default-roles-zpi-realm"`, `"realm-zpi-role"`. Jeśli przynajmniej jedna z wczytanych ról daje użytkownikowi dostęp do określonego zapytania, to użytkownik jest pomyślnie autoryzowany. W przeciwnym wypadku zwracany jest błąd 403 - jeżeli żadna z ról nie nadaje mu odpowiednich uprawnień.

### Możliwe uprawnienia w aplikacji
W aplikacji wyróżniamy akcje create, read, update, delete, list. Każde z uprawnień akcji musi być nadane osobno, w szczególności uprawnienie do listowania zasobów nie daje automatycznie możliwości przeglądania ich szczegółów. Uprawnienia można definiować dla konkretnej przestrzeni nazw i typu zasobu, np. Pod, ConfigMap.

### Mapa ról na uprawnienia
W celu korzystania z aplikacji należy przed instalacją zdefiniować zasób typu ConfigMap z kluczami role-map i subrole-map (opcjonalnie). Mapę ról można zmieniać w trakcie działania aplikacji. Nazwę i namespace ConfigMapy należy zdefiniować w zmiennych środowiskowych `ROLEMAP_NAMESPACE` i `ROLEMAP_NAME`. Domyślne wartości to namespace "default" i nazwa "role-map". Przykładowa definicja:
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
Definiując operacje, można uwzględnić trzy atrybuty: namespace, resource (typ zasobu, np. "Pod") oraz akcje (CRUDL). Każdy z tych atrybutów można pominąć, co będzie jednoznaczne z nadaniem/ograniczeniem uprawnień dla wszystkich możliwych namespace'ów, resource'ów lub akcji. Taki sam efekt można uzyskać, wpisując "*" dla namespace i resource lub ["*"] dla operations. Przykład:
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
Według powyższej definicji roli, `admin` ma prawo do wszystkich akcji na wszystkich zasobach w dowolnym namespace, za wyjątkiem namespace `"top-restricted"` oraz usuwania, edytowania lub tworzenia ConfigMap w namespace `"role-map-namespace"`. Wartość `operations: ["*"]` jest wymagana, gdyż potrzebny jest przynajmniej jeden atrybut z namespace, resource, operations.

### Używanie podról
Używając podról, można zdefiniować konfiguracje uprawnień często powtarzających się pomiędzy poszczególnymi rolami. Ważne jest rozróżnienie pomiędzy rolą a podrolą - nazwa roli pochodzi z zewnętrznego dostawcy tożsamości i musi być dokładnie taka sama jak w tokenie JWT, aby gwarantować użytkownikowi jakiekolwiek uprawnienia. Podrola służy wyłącznie do przekazywania uprawnień do roli. Można zdefiniować rolę i podrolę o tej samej nazwie. Aby rola otrzymała uprawnienia z podroli, należy dodać nazwę podroli do listy "subroles" w roli. Nie można używać ról jako podról. Podrole też mogą posiadać podrole.

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
W przedstawionym przykładzie `team1admin` jest jednocześnie zdefiniowany jako rola i podrola. Nie powoduje to problemu, gdyż role i podrole są zawsze rozpatrywane jako osobne byty. Podrola `team1admin` dziedziczy uprawnienia z `permissionViewer`, umożliwiając dowolne operacje w namespace `team1` oraz odczyt i listowanie ConfigMap w namespace `"role-map-namespace"`. Analogicznie działa `team2admin`, operując w namespace `"team2"`. Podrole te są podrolami dla ról `team1admin`, `team2admin` oraz dla `manager`.

Rola `manager` niezależnie od uprawnień podról (`team1admin`, `team2admin`) nie może wykonywać operacji delete, create ani update, co wynika z nadrzędnych ograniczeń (deny). Ograniczenia te mają wyższy priorytet niż uprawnienia z podról. Podrole mogą dziedziczyć uprawnienia od innych podról, co widać na przykładzie `team1Admin` i `permissionViewer`.

Ostatecznie rola `team1admin` gwarantuje uprawnienia do wszystkich zasobów i akcji w namespace `"team1"` oraz podgląd i listowanie ConfigMap w namespace `"role-map-namespace"`, analogicznie `team2admin` dla namespace `"team2"`. Rola `manager` może listować i czytać szczegóły zasobów w namespace’ach `"team1"`, `"team2"` oraz ConfigMap’y w namespace `"role-map-namespace"`

## Authorization in the Application

### JWT Token

The application connects to an identity provider such as Keycloak or another provider compliant with the OpenID Connect protocol that provides user roles. Authorization is based on extracting user roles from the JWT token sent with every request to the backend server. The application reads roles defined in two sections of the token:

- **`resource_access` → `VITE_KEYCLOAK_CLIENTNAME` → `roles`**: A list of roles assigned in the context of a specific client.
- **`realm_access` → `roles`**: A list of roles assigned globally within the realm.

The value of `VITE_KEYCLOAK_CLIENTNAME` corresponds to the client name configured as an environment variable in the Helm Chart. This variable indicates the token section relevant to the application.
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
In this case, for `VITE_KEYCLOAK_CLIENTNAME` equal to `"ZPI-client"`, the application would authorize the request based on the roles `"zpi-role"`, `"default-roles-zpi-realm"`, and `"realm-zpi-role"`. If at least one of the extracted roles grants the user access to the requested operation, the user is successfully authorized. Otherwise, a `403 Forbidden` error is returned if none of the roles provide the necessary permissions.

### Possible Permissions in the Application

The application distinguishes the following actions: **create**, **read**, **update**, **delete**, and **list**. Each action's permission must be granted separately. Notably, permission to list resources does not automatically grant access to view their details. Permissions can be defined for specific namespaces and resource types, such as `Pod` or `ConfigMap`.

### Role-to-Permissions Mapping

To use the application, you must define a `ConfigMap` resource with the keys `role-map` and optionally `subrole-map` before installation. The role map can be modified while the application is running. The name and namespace of the `ConfigMap` must be set in the environment variables `ROLEMAP_NAMESPACE` and `ROLEMAP_NAME`. The default values are the namespace `"default"` and the name `"role-map"`. Example definition:
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
The fields for defining roles and subroles are the same. These are:

- **Role/Subrole Name** - Provided as a key before permissions. For roles, it must match exactly the role name extracted from the list of roles in the JWT token.
- **permit** - A list of operations that the role is allowed to perform.
- **deny** - A list of operations that the role is prohibited from performing.
- **subroles** - A list of subrole names from which permissions are inherited.

In addition to the role/subrole name, at least one attribute (`permit`, `deny`, or `subroles`) must be defined for the role.

### Defining Operations in `permit` and `deny`

When defining operations, you can specify three attributes: **namespace**, **resource** (resource type, e.g., `Pod`), and **action** (CRUDL). Any of these attributes can be omitted, which will grant or restrict permissions for all possible namespaces, resources, or actions. The same effect can be achieved by using `*` for namespace and resource or `[“*”]` for operations. Example:
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
According to the above role definition, the `admin` role has permission to perform all actions on all resources in any namespace, except for the namespace `top-restricted` and for deleting, editing, or creating `ConfigMap` resources in the namespace `role-map-namespace`. The value `operations: ["*"]` is required, as at least one of the attributes `namespace`, `resource`, or `operations` must be specified.

### Using Subroles

Subroles can be used to define permission configurations that are frequently repeated across different roles. It is important to distinguish between a role and a subrole. The role name is provided by an external identity provider and must match exactly the name in the JWT token to grant the user any permissions. Subroles are used solely for passing permissions to a role. It is possible to define a role and a subrole with the same name. For a role to inherit permissions from a subrole, the subrole's name must be added to the `subroles` list in the role. Roles cannot be used as subroles. Subroles, however, can have their own subroles.

Restrictions (`deny`) defined in a subrole do not affect permissions (`permit`) defined in the parent role or permissions inherited from other subroles of the parent role. However, restrictions (`deny`) from the parent role do affect permissions (`permit`) inherited from all its subroles.

Example:
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
In the example above, `team1Admin` is defined as both a role and a subrole. This does not cause a problem, as roles and subroles are always considered separate entities. The subrole `team1admin` inherits permissions from `permissionViewer`, allowing any operations in the `team1` namespace and reading and listing `ConfigMap` resources in the `role-map-namespace`. The same applies to `team2admin`, operating in the `team2` namespace. These subroles are subroles for the roles `team1admin`, `team2admin`, and `manager`.

The `manager` role, regardless of the permissions of its subroles (`team1admin`, `team2admin`), cannot perform delete, create, or update actions. This is due to the parent role's restrictions (deny), which take precedence over subrole permissions. Subroles can inherit permissions from other subroles, as seen in the `team1Admin` and `permissionViewer` example.

Ultimately, the `team1Admin` role grants permissions to all resources and actions in the `team1` namespace and viewing and listing `ConfigMap` resources in the `role-map-namespace`. Similarly, `team2admin` grants permissions in the `team2` namespace. The `manager` role can list and read resource details in the `team1` and `team2` namespaces and `ConfigMap` resources in the `role-map-namespace`.