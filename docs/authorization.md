## For english press [here](#authorization-in-kubernetes-access-manager)
## Opis autoryzacji w Kubernetes Access Manager
### Token JWT

KAM zakłada połączenie z dostarczycielem tożsamości Keycloak lub innym spełniającym wymagania protokołu OpenIdConnect i dostarczającego role. Autoryzacja oparta jest na ekstrakcji ról użytkownika z tokena JWT wysyłanego z każdym zapytaniem do serwera backend. KAM odczytuje role zdefiniowane w dwóch miejscach tokena:

- resource_access → `VITE_KEYCLOAK_CLIENTNAME` → roles: Lista ról przypisanych w kontekście konkretnego klienta.
- realm_access → roles: Lista ról przypisanych globalnie w ramach danego realm.

Zmienna środowiskowa `VITE_KEYCLOAK_CLIENTNAME` odpowiada nazwie klienta dodanego w Keycloak dla KAM. Zmienna ta wskazuje na sekcję tokena odpowiadającą KAM. Przykład:
```json
{
  ...
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
  ...
}
```
W tym wypadku dla `VITE_KEYCLOAK_CLIENTNAME` równego `ZPI-client` KAM autoryzowałby żądanie na podstawie ról `zpi-role`, `default-roles-zpi-realm`, `realm-zpi-role`. Jeśli przynajmniej jedna z wczytanych ról daje użytkownikowi dostęp do określonego zapytania, to użytkownik jest pomyślnie autoryzowany. W przeciwnym wypadku zwracany jest błąd zgodnie z definicją [API](./api-swagger.yaml).

### Możliwe uprawnienia w KAM
W KAM wyróżniamy akcje create, read, update, delete, list. Każde z uprawnień akcji musi być nadane osobno, w szczególności uprawnienie do listowania zasobów nie daje automatycznie możliwości przeglądania ich szczegółów. Uprawnienia można definiować dla konkretnej przestrzeni nazw i typu zasobu, np. Pod, ConfigMap.

### Mapa ról na uprawnienia
W celu korzystania z KAM należy przed instalacją zdefiniować zasób typu ConfigMap z kluczami role-map i subrole-map (opcjonalnie). Konfigurację można zmieniać w trakcie działania KAM. Nazwę i namespace ConfigMapy należy zdefiniować w zmiennych środowiskowych `ROLEMAP_NAMESPACE` i `ROLEMAP_NAME`. Domyślne wartości to namespace "default" i nazwa "role-map".

Pola definicji roli i podroli są takie same. Są to:
- **nazwa roli/podroli** - podana przed uprawnieniami, jako klucz. W przypadku roli powinna być dokładnie taka sama jak rola odczytana z listy ról w tokenie JWT.
- **permit** - lista operacji na które zezwala dana rola
- **deny** - lista operacji które rola zabrania
- **subroles** - lista nazw podról, z których dziedziczone są uprawnienia.

Oprócz samej nazwy roli/podroli należy zdefiniować przynajmniej jeden atrybut (permit/deny/subroles). Przykład definicji jednej roli/podroli:
```yaml
    admin:          # nazwa roli
      permit:       # lista operacji dozwolonych
        - namespace: "namespace"
          resource: "*"
          operations: ["*"]
        - namespace: "namespace2"
          resource: "Pod"
          operations: ["read", "list"]
      deny:         # lista operacji zabronionych
        - namespace: "namespace"
          resource: "ConfigMap"
          operations: ["delete", "create", "update"]
```

### Definiowanie operacji w `permit`, `deny`
Definiując operacje, można uwzględnić trzy atrybuty: namespace, resource (typ zasobu, np. "Pod") oraz akcje ("create", "read", "update", "delete", "list"). Każdy z tych atrybutów można pominąć, co będzie jednoznaczne z nadaniem/ograniczeniem uprawnień dla wszystkich możliwych namespace'ów, resource'ów lub akcji. Taki sam efekt można uzyskać, wpisując `*` dla namespace i resource lub `["*"]` dla operations. Przykład:
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
Według powyższej definicji roli, `admin` ma prawo do wszystkich akcji na wszystkich zasobach w dowolnym namespace, za wyjątkiem namespace `top-restricted` oraz usuwania, edytowania lub tworzenia ConfigMap w namespace `role-map-namespace`. Wartość `operations: ["*"]` jest wymagana, gdyż potrzebny jest przynajmniej jeden atrybut z namespace, resource, operations.

### Używanie podról

Używając podról, można zdefiniować konfiguracje uprawnień, które są często powtarzane pomiędzy poszczególnymi rolami. Ważne jest rozróżnienie pomiędzy rolą a podrolą: nazwa roli pochodzi od zewnętrznego dostawcy tożsamości i musi być dokładnie taka sama jak w tokenie JWT, aby użytkownik mógł uzyskać jakiekolwiek uprawnienia. Podrola natomiast służy wyłącznie do przekazywania uprawnień do roli. Można zdefiniować zarówno rolę, jak i podrolę o tej samej nazwie. Aby rola otrzymała uprawnienia z podroli, należy dodać nazwę tej podroli do listy `subroles` w konfiguracji roli. Nie można używać ról jako podról. Podrole mogą posiadać własne podrole.

Ograniczenia `deny` zdefiniowane w podroli nie wpływają na uprawnienia `permit` zdefiniowane w nadrzędnej roli ani na uprawnienia wynikające z innych podról przypisanych do tej roli. Natomiast ograniczenia `deny` w nadrzędnej roli mają wpływ na wszystkie uprawnienia `permit` zdefiniowane w jej podrolach. 

Przykłady:

#### 1. Prosty przypadek użycia podroli
```yaml
  role-map: |
    userWithList:
      permit:
        - operations: ["list"]
      subroles:
        - "permissionsViewer"
    user:
      subroles:
        - "permissionsViewer"
  subrole-map: |
    permissionsViewer:
      permit:
        - namespace: "role-map-namespace"
          resource: "ConfigMap"
          operations: ["read", "list"]
```
W powyższym przykładzie `user` i `userWithList` otrzymują uprawnienia zdefiniowane w podroli `permissionsViewer` - czytanie i listowanie `ConfigMap` w namespace `role-map-namespace`. Rola `userWithList` dodatkowo ma uprawnienie do listowania wszystkich zasobów.

#### 2. Ograniczanie uprawnień z podról
```yaml
  role-map: |
    role:
      permit:
        - operations: ["list"]
      deny:
        - namespace: "other-restricted"
          operations: ["*"]
      subroles:
        - "readCreator"
  subrole-map: |
    readCreator:
      deny:
        - namespace: "restricted"
          operations: ["*"]
      permit:
        - operations: ["read", "create"]
```
W powyższym przykładzie podrola `readCreator` nadaje uprawnienia do akcji `read` i `create` w dowolnym namespace poza `restricted`. Rola `role` posiada swoje własne uprawnienia do listowania i ograniczenie dostępu do namespace'a `other-restricted` oraz otrzymuje uprawnienia z podroli `readCreator`. Oznacza to że użytkownik z rolą `role` będzie mógł:
- listować zasoby w dowolnym namespace poza `other-restricted` zgodnie z zasadą, że ograniczenia z podroli nie mają wpływu na uprawnienia z nadroli
- czytać i tworzyć zasoby w dowolnym namespace poza `restricted` i `other-restricted` zgodnie z zasadą, że ograniczenia z nadroli mają wpływ na uprawnienia z podroli

#### 3. Rozróżnienie pomiędzy rolą a podrolą, podrole podról
Przykład:
```yaml
  role-map: |
    manager:
      deny:
        - operations: ["delete", "create", "update"]
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
W przedstawionym przykładzie `team1admin` jest jednocześnie zdefiniowany jako rola i podrola. Nie powoduje to problemu, gdyż role i podrole są zawsze rozpatrywane jako osobne byty. Podrola `team1admin` dziedziczy uprawnienia z `permissionViewer`, umożliwiając dowolne operacje w namespace `team1` oraz odczyt i listowanie ConfigMap w namespace `role-map-namespace`. Analogicznie działa `team2admin`, operując w namespace `team2`. Podrole te są podrolami dla ról `team1admin`, `team2admin` oraz dla `manager`.

Rola `manager` niezależnie od uprawnień podról (`team1admin`, `team2admin`) nie może wykonywać operacji delete, create ani update, co wynika z nadrzędnych ograniczeń (deny). Ograniczenia te mają wyższy priorytet niż uprawnienia z podról. Podrole mogą dziedziczyć uprawnienia od innych podról, co widać na przykładzie `team1Admin` i `permissionViewer`.

Ostatecznie rola `team1admin` gwarantuje uprawnienia do wszystkich zasobów i akcji w namespace `team1` oraz podgląd i listowanie ConfigMap w namespace `role-map-namespace`, analogicznie `team2admin` dla namespace `team2`. Rola `manager` może listować i czytać szczegóły zasobów w namespace’ach `team1`, `team2` oraz ConfigMap’y w namespace `role-map-namespace`

### Pełne wykorzystanie wszystkich możliwości definiowanie mapy ról
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

## Authorization in Kubernetes Access Manager
### JWT token

KAM assumes a connection with an identity provider such as Keycloak or another provider that complies with the OpenIdConnect protocol and provides roles. Authorization is based on extracting user roles from the JWT token sent with each request to the backend server. KAM reads roles defined in two sections of the token:

- `resource_access` → `VITE_KEYCLOAK_CLIENTNAME` → roles: A list of roles assigned in the context of a specific client.
- `realm_access` → roles: A list of roles assigned globally within a given realm.

The environment variable `VITE_KEYCLOAK_CLIENTNAME` corresponds to the name of the client added in Keycloak for KAM. This variable points to the token section associated with KAM. Example:

```json
{
  ...
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
  ...
}
```
In this case, for `VITE_KEYCLOAK_CLIENTNAME` equal to `ZPI-client`, KAM would authorize the request based on the roles `zpi-role`, `default-roles-zpi-realm`, and `realm-zpi-role`. If at least one of the loaded roles grants the user access to a specific request, the user is successfully authorized. Otherwise, an error is returned according to the [API definition](./api-swagger.yaml).

### Possible permissions in KAM

KAM distinguishes between the following actions: create, read, update, delete, and list. Each action's permission must be granted separately. In particular, having permission to list resources does not automatically grant the ability to view their details. Permissions can be defined for a specific namespace and resource type, e.g., Pod, ConfigMap.

### Mapping roles to permissions

To use KAM, a resource of type ConfigMap must be defined before installation, with the keys `role-map` and optionally `subrole-map`. The configuration can be modified while KAM is running. The name and namespace of the ConfigMap must be specified in the environment variables `ROLEMAP_NAMESPACE` and `ROLEMAP_NAME`. Default values are the namespace `default` and the name `role-map`.

The fields for defining roles and subroles are the same and include the following:

- **Role/Subrole Name** - Specified before the permissions, as a key. For roles, it must match exactly the name of the role read from the list of roles in the JWT token.
- **permit** - A list of operations that the role allows.
- **deny** - A list of operations that the role prohibits.
- **subroles** - A list of subrole names from which permissions are inherited.

In addition to the role/subrole name, at least one attribute (`permit`, `deny`, or `subroles`) must be defined. An example definition of a single role/subrole:

```yaml
    admin:          # role name
      permit:       # permitted operations list
        - namespace: "namespace"
          resource: "*"
          operations: ["*"]
        - namespace: "namespace2"
          resource: "Pod"
          operations: ["read", "list"]
      deny:         # denied operations list
        - namespace: "namespace"
          resource: "ConfigMap"
          operations: ["delete", "create", "update"]
```

### Defining operations in `permit` and `deny`

When defining operations, three attributes can be specified: `namespace`, `resource` (the type of resource, e.g., "Pod"), and `operations` ("create", "read", "update", "delete", "list"). Each of these attributes can be omitted, which will be interpreted as granting or restricting permissions for all possible namespaces, resources, or actions. The same effect can be achieved by explicitly using `*` for `namespace` and `resource`, or `["*"]` for `operations`. 

Example:

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
According to the above role definition, `admin` has the right to perform all actions on all resources in any namespace, except for the namespace `top-restricted` and the actions of deleting, editing, or creating ConfigMaps in the namespace `role-map-namespace`. The value `operations: ["*"]` is required, as at least one attribute from `namespace`, `resource`, or `operations` must be specified.

### Using subroles

By using subroles, you can define configurations of permissions that are frequently reused across various roles. It is important to distinguish between a role and a subrole: the role name is derived from an external identity provider and must match exactly the role name in the JWT token for the user to gain any permissions. A subrole, on the other hand, is used solely to pass permissions to a role. Both a role and a subrole can be defined with the same name. To grant a role permissions from a subrole, the name of the subrole must be added to the `subroles` list in the role configuration. Roles cannot be used as subroles, but subroles can have their own subroles.

The `deny` restrictions defined in a subrole do not affect the `permit` permissions defined in the parent role or those resulting from other subroles assigned to that role. However, `deny` restrictions in the parent role affect all `permit` permissions defined in its subroles.

Examples:

#### 1. Simple use case of subroles
```yaml
  role-map: |
    userWithList:
      permit:
        - operations: ["list"]
      subroles:
        - "permissionsViewer"
    user:
      subroles:
        - "permissionsViewer"
  subrole-map: |
    permissionsViewer:
      permit:
        - namespace: "role-map-namespace"
          resource: "ConfigMap"
          operations: ["read", "list"]
```
In the example above, both `user` and `userWithList` inherit permissions defined in the subrole `permissionsViewer` reading and listing `ConfigMap` resources in the `role-map-namespace` namespace. Additionally, the `userWithList` role has the permission to list all resources.

#### 2. Restricting permissions from subroles
```yaml
  role-map: |
    role:
      permit:
        - operations: ["list"]
      deny:
        - namespace: "other-restricted"
          operations: ["*"]
      subroles:
        - "readCreator"
  subrole-map: |
    readCreator:
      deny:
        - namespace: "restricted"
          operations: ["*"]
      permit:
        - operations: ["read", "create"]
```
In the example above, the subrole `readCreator` grants permissions for the `read` and `create` actions in any namespace except `restricted`. The `role` role has its own permissions to list resources and restricts access to the `other-restricted` namespace. Additionally, it inherits permissions from the `readCreator` subrole. This means that a user with the `role` role will be able to:

- List resources in any namespace except `other-restricted`, according to the rule that restrictions from subroles do not affect the permissions granted by the parent role.
- Read and create resources in any namespace except `restricted` and `other-restricted`, following the rule that restrictions in the parent role affect the permissions inherited from subroles.

#### 3. Distinction between roles and subroles, subroles of subroles
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
In the presented example, `team1admin` is simultaneously defined as both a role and a subrole. This does not cause any issues, as roles and subroles are always treated as separate entities. The subrole `team1admin` inherits permissions from `permissionViewer`, allowing all operations in the `team1` namespace as well as read and list access to ConfigMaps in the `role-map-namespace` namespace. Similarly, `team2admin` operates in the `team2` namespace. These subroles are used as subroles for the `team1admin`, `team2admin`, and `manager` roles.

The `manager` role, regardless of the permissions of its subroles (`team1admin`, `team2admin`), cannot perform delete, create, or update operations due to the overarching restrictions (`deny`). These restrictions take precedence over permissions inherited from subroles. Subroles can inherit permissions from other subroles, as demonstrated in the relationship between `team1admin` and `permissionViewer`.

Ultimately, the `team1admin` role grants permissions for all resources and actions in the `team1` namespace and provides read and list access to ConfigMaps in the `role-map-namespace` namespace. Similarly, `team2admin` provides equivalent permissions for the `team2` namespace. The `manager` role can list and read the details of resources in the `team1` and `team2` namespaces, as well as ConfigMaps in the `role-map-namespace`.

### Fully leveraging the possibilities of defining the role map

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