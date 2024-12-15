## For English Press [Here]()

# Architektura

## Ogólna struktura projektu
![Diagram komponentów](images/component_diagram.png)

Kubernetes Access Manager (KAM) składa się z trzech głównych komponentów:
- **Frontend**- aplikacja webowa, która pozwala na zarządzanie zasobami w klastrze Kubernetes. 
Frontend jest odpowiedzialny za wyświetlanie interfejsu użytkownika, a także za komunikację z backendem.
Do uwierzytelniania użytkowników wykorzystuje on tokeny JWT.
- **Backend**- serwer aplikacyjny, który zarządza zasobami w klastrze Kubernetes.
Backend jest odpowiedzialny za autoryzację użytkowników, a także za komunikację z API klastra Kubernetes.
Do uwierzytelniania użytkowników wykorzystuje on tokeny JWT.
- **Identity Provider**- dostawca tożsamości, który pozwala na uwierzytelnianie i autoryzację użytkowników.
Wybrany został protokół OpenID Connect, a jako dostawca tożsamości wykorzystany został Keycloak.

Kam komunikuje się z API klastra Kubernetes za pomocą biblioteki `client-go`, natomiast z Klientem Helmowym za pomocą biblioteki `helm`.

## Helm chart
Aplikacja uruchamiana jest w klastrze Kubernetes. W celu łatwej instalacji i zarządzania aplikacją wykorzystany został Helm.

![Diagram Helm Chart](images/helm_chart.png)

Helm chart aplikacji składa się z czterech głównych komponentów:
- **Backend**- serwer aplikacyjny, który zarządza zasobami w klastrze Kubernetes. Składają się na niego:
    - Pod- w nich uruchamiana jest instancja serwera backendowego.
    - Deployment- definiuje, w jaki sposób mają być tworzone i zarządzane instancje podów.
    - Service- umożliwia komunikację pomiędzy częścią frontendową a backendową wewnątrz klastra Kubernetes.
    - Horizontal Pod Autoscaler- umożliwia automatyczne skalowanie liczby podów w zależności od obciążenia.
- **Frontend**- serwer aplikacji webowej, która pozwala na zarządzanie zasobami w klastrze Kubernetes. Składają się na niego:
    - Pod- w nich uruchamiana jest instancja serwera frontendowego.
    - Deployment- definiuje, w jaki sposób mają być tworzone i zarządzane instancje podów.
    - Service- umożliwia komunikację pomiędzy częścią frontendową a backendową wewnątrz klastra Kubernetes.
    - Horizontal Pod Autoscaler- umożliwia automatyczne skalowanie liczby podów w zależności od obciążenia.
- **Access Control**- zbiór zasobów przyznający uprawnienia wymagane w części backendowej. Składają się na niego:
    - Service Account-umożliwia aplikacji dostęp do API klastra Kubernetes.
    - Cluster Role- definiuje zestaw uprawnień wymaganych do zarządzania zasobami w klastrze Kubernetes.
    - Cluster Role Binding- definiuje powiązanie pomiędzy Service Accountem a Cluster Role.
- **Ingress**- definiuje reguły zarządzania ruchem przychodzącym do klastra Kubernetes.

Część **Access Control** jest wymagana tylko w przypadku, gdy na klastrze włączony jest mechanizm RBAC.

## Część backendowa
REST API napisane w języku Go, które pozwala na zarządzanie zasobami w klastrze Kubernetes.
Może funkcjonować w dwóch trybach:
- **Wewnątrz klastra**- domyślny tryb w obrazach dockerowych, aplikacja pobiera konfigurację bezpośrednio z klastra Kubernetes.
- **Zewnętrzny**- domyślny tryb przy zwykłym uruchamianiu, przydatny podczas rozwijania aplikacji lokalnie. Konfiguracja klastra pobierana jest domyślnie z pliku `kubeconfig`.
Możliwe jest również podanie ścieżki do pliku.

Backend komunikuje się z API klastra Kubernetes za pomocą biblioteki `client-go`, natomiast z Klientem Helmowym za pomocą biblioteki `helm`.

### Konfiguracja klastra Kubernetes
Konfiguracja klastra Kubernetes przechowywana jest w singletonie `ClientSingleton`. Sposób pobierania konfiguracji zależy od flagi `--in-cluster`. 
Domyślnie wartość flagi wynosi `false`, co oznacza, że konfiguracja pobierana jest z pliku `kubeconfig`.
W przypadku, gdy flaga przyjmuje wartość `true`, konfiguracja pobierana jest bezpośrednio z klastra Kubernetes.

Lokalizacja pliku `kubeconfig` zależy od dwóch wartości:
- Flagi `--kubeconfig`- ścieżka do pliku `kubeconfig`.
- Zmiennej środowiskowej `KUBECONFIG`- ścieżka do pliku `kubeconfig`.

Jeżeli żadna z tych wartości nie jest podana, domyślnie aplikacja korzysta z pliku `~/.kube/config`.

### Funkcje create, get, update, delete klastra Kubernetes
Funkcje te wykorzystują bibliotekę `client-go` do komunikacji z API klastra Kubernetes. Jako argumenty przyjmują kolejno:
- GetResource- typ zasobu (resourceType), namespace, nazwę zasobu (resourceName).
- CreateResource- typ zasobu (resourceType), namespace, zasób (ResourceDetails).
- UpdateResource- typ zasobu (resourceType), namespace, nazwę zasobu (resourceName), zasób (ResourceDetails).
- DeleteResource- typ zasobu (resourceType), namespace, nazwę zasobu (resourceName).

W przypadku gdy zasób jest namespace'owany, a użytkownik nie poda namespace, zasób zostanie utworzony w namespace'u `default`.

ResourceDetails to luźna struktura, która reprezentuje nieustrukturyzowane dane zasobu.

Funkcje create, get i update zwracają zasób w postaci struktury `ResourceDetails` lub błąd w przypadku niepowodzenia. Funkcja delete zwraca błąd w przypadku niepowodzenia.

Funkcje te są w stanie obsłużyć wszystkie typy zasobów, jednak ze względu na funkcje listowania zakres został ograniczony do 20 typów wymienionych w części `Funkcja list`.

### Funkcja list klastra Kubernetes
Funkcja ta wykorzystuje bibliotekę `client-go` do komunikacji z API klastra Kubernetes. 
Jako argumenty przyjmuje: typ zasobu (resourceType), namespace.

W przypadku gdy zasób jest namespace'owany, a użytkownik nie poda namespace, zwracane są zasoby ze wszystkich namespace'ów.

Funkcja dla każdego typu zasobu zwraca określoną listę wartości oraz listę nazw zwracanych pól.
- Dla zasobów typu `ReplicaSet` zwracane są wartości `name`, `namespace`, `desired`, `current`, `ready`, `age`.
- Dla zasobów typu `Pod` zwracane są wartości `name`, `namespace`, `containers`, `restarts`, `controlled_by`, `node`, `qos`, `age`, `status`.
- Dla zasobów typu `Deployment` zwracane są wartości `name`, `namespace`, `pods`, `replicas`, `age`, `conditions`.
- Dla zasobów typu `ConfigMap` zwracane są wartości `name`, `namespace`, `keys`, `age`.
- Dla zasobów typu `Secret` zwracane są wartości `name`, `namespace`, `labels`, `keys`, `type`, `age`.
- Dla zasobów typu `Ingress` zwracane są wartości `name`, `namespace`, `loadbalancers`, `age`.
- Dla zasobów typu `PersistentVolumeClaim` zwracane są wartości `name`, `namespace`, `storage_class`, `size`, `age`, `status`.
- Dla zasobów typu `StatefulSet` zwracane są wartości `name`, `namespace`, `pods`, `replicas`, `age`.
- Dla zasobów typu `DaemonSet` zwracane są wartości `name`, `namespace`, `pods`, `node_selector`, `age`.
- Dla zasobów typu `Job` zwracane są wartości `name`, `namespace`, `completions`, `age`, `conditions`.
- Dla zasobów typu `CronJob` zwracane są wartości `name`, `namespace`, `schedule`, `suspend`, `active`, `last_schedule`, `age`.
- Dla zasobów typu `Service` zwracane są wartości `name`, `namespace`, `type`, `cluster_ip`, `ports`, `external_ip`, `selector`, `age`.
- Dla zasobów typu `ServiceAccount` zwracane są wartości `name`, `namespace`, `age`.
- Dla zasobów typu `Node` zwracane są wartości `name`, `taints`, `roles`, `version`, `age`, `conditions`.
- Dla zasobów typu `Namespace` zwracane są wartości `name`, `labels`, `status`, `age`.
- Dla zasobów typu `CustomResourceDefinition` zwracane są wartości `resource`, `group`, `version`, `scope`, `age`.
- Dla zasobów typu `PersistentVolume` zwracane są wartości `name`, `storage_class`, `capacity`, `claim`, `age`, `status`.
- Dla zasobów typu `StorageClass` zwracane są wartości `name`, `provisioner`, `reclaim_policy`, `default`, `age`.
- Dla zasobów typu `ClusterRole` zwracane są wartości `name`, `age`.
- Dla zasobów typu `ClusterRoleBinding` zwracane są wartości `name`, `bindings`, `age`.

### Funkcje pomocnicze klastra Kubernetes
- **GetResourceGroupVersion**- na podstawie typu zasobu zwraca informacje potrzebne do generycznego wywołania funkcji `client-go`.
Dodatkowo zwraca informacje o tym, czy zasób jest namespace'owany, czy nie.
- **GetResourceInterface**- dodatkowy poziom abstrakcji umożliwiający wstrzykiwanie zależności (dependency injection) w testach jednostkowych.
Równocześnie upraszcza i ujednolica wywołanie funkcji `client-go`.
- **WatchForChanges**- funkcja umożliwiająca obserwowanie zmian w zasobie opisującym role wykorzystywane w autoryzacji.
- **transposeResourceListColumns**- funkcja pomocnicza, która zamienia listę zasobów i ich pól (postaci mapy nazwa_zasobu -> jej_pola) na listę pól i do jakich typów zasobów należą (postaci mapy nazwa_pola -> lista_typów).

### Konfiguracja Akcji Helm
Konfiguracja ta twożona jest za pomocą funkcji `getActionConfig` na podstawie konfiguracji klastra Kubernetes oraz namespace'a, w którym ma zostać wykonana akcja.

### Funkcje get, rollback, uninstall, getHistory Helm
Funkcje te wykorzystują bibliotekę `helm` do komunikacji z Klientem Helmowym. Jako argumenty przyjmują kolejno:
- GetRelease- nazwę releasu (releaseName).
- RollbackRelease- nazwę releasu (releaseName), numer wersji (revision).
- UninstallRelease- nazwę releasu (releaseName).
- GetReleaseHistory- nazwę releasu (releaseName), ilość wersji (max).

