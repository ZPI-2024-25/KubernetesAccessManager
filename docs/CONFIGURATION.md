## For English press [here](#configuration)
# Konfiguracja

Ten dokument zawiera przegląd wartości konfigurowalnych w pliku [`values.yaml`](../charts/kam/values.yaml) dla chartu Helm. Te wartości mogą być dostosowywane w celu skonfigurowania wdrożenia aplikacji. W celu uzyskania szczegółowych informacji jak wdrożyć aplikację, odwołaj się do pliku [DEPLOYMENT.md](./DEPLOYMENT.md).

## Ogólna konfiguracja

- **nameOverride**: Nadpisuje nazwę chartu. Przydatne, jeśli chcesz wdrożyć wiele instancji chartu z różnymi nazwami.
- **fullnameOverride**: Całkowicie nadpisuje wygenerowaną nazwę chartu. Przydatne do zapewnienia spójnych konwencji nazewnictwa.

## Globalna konfiguracja

### **global.env.KEYCLOAK_JWKS_URL**
- **Opis**: URL do zestawu JSON Web Key Set (JWKS) w Keycloak. Ten URL jest używany do pobierania kluczy publicznych w celu weryfikacji tokenów JWT wydanych przez Keycloak.
- **Wymagane**: Jeśli `VITE_KEYCLOAK_URL`, `VITE_KEYCLOAK_REALMNAME` i `VITE_KEYCLOAK_CLIENTNAME` nie są podane.
- **Domyślne**: Generowane na podstawie `VITE_KEYCLOAK_URL`, `VITE_KEYCLOAK_REALMNAME` i `VITE_KEYCLOAK_CLIENTNAME`.
- **Używane przez**: Backend
- **Przykład**: `https://keycloak.example.com/realms/myrealm/protocol/openid-connect/certs`

### **global.env.VITE_KEYCLOAK_URL**
- **Opis**: Bazowy URL dla serwera Keycloak. Jest używany przez backend do interakcji z Keycloak w celu uwierzytelniania i autoryzacji. Powinien być dostępny z backendu.
- **Wymagane**: Jeśli `KEYCLOAK_JWKS_URL` nie jest podany.
- **Domyślne**: Brak
- **Używane przez**: Backend
- **Przykład**: `https://keycloak.example.com`

### **global.env.FRONTEND_KEYCLOAK_URL**
- **Opis**: URL do serwera Keycloak używanego przez frontend. Jest wykorzystywany przez frontend do celów uwierzytelniania. Powinien być publicznie dostępny.
- **Wymagane**: Tak
- **Domyślne**: Brak
- **Używane przez**: Frontend
- **Przykład**: `https://keycloak.example.com`

### **global.env.VITE_KEYCLOAK_REALMNAME**
- **Opis**: Nazwa realm w Keycloak. Realm w Keycloak to przestrzeń, w której zarządzasz obiektami, takimi jak użytkownicy, role i klienci.
- **Wymagane**: Tak
- **Domyślne**: Brak
- **Używane przez**: Frontend i Backend
- **Przykład**: `myrealm`

### **global.env.VITE_KEYCLOAK_CLIENTNAME**
- **Opis**: Nazwa klienta w Keycloak. Klient w Keycloak to podmiot, który może poprosić o uwierzytelnienie użytkownika. Należy utworzyć klienta w Keycloak dla naszej aplikacji.
- **Wymagane**: Tak
- **Domyślne**: Brak
- **Używane przez**: Frontend i Backend
- **Przykład**: `myclient`

### **global.env.ROLEMAP_NAMESPACE**
- **Opis**: Namespace, w którym jest przechowywana mapa ról. Służy do określenia namespace w Kubernetes, gdzie znajduje się ConfigMap z mapą ról.
- **Wymagane**: Nie
- **Domyślne**: `default`
- **Używane przez**: Backend
- **Przykład**: `mynamespace`

### **global.env.ROLEMAP_NAME**
- **Opis**: Nazwa mapy ról. Służy do określenia nazwy ConfigMap, która zawiera mapę ról.
- **Wymagane**: Nie
- **Domyślne**: `role-map`
- **Używane przez**: Backend
- **Przykład**: `myrolemap`

## Konfiguracja Backend

- **backend.replicaCount**: Liczba replik dla wdrożenia backendu.
- **backend.image.repository**: Repozytorium obrazu Docker dla backendu.
- **backend.image.pullPolicy**: Polityka pobierania obrazu dla backendu. Opcje: `Always`, `IfNotPresent`, `Never`.
- **backend.image.tag**: Tag obrazu Docker używanego dla backendu. Domyślnie: wersja aplikacji w chart.
- **backend.serviceAccount.create**: Czy utworzyć nowe konto serwisowe dla backendu.
- **backend.serviceAccount.automount**: Czy automatycznie montować API credentials konta serwisowego.
- **backend.serviceAccount.annotations**: Adnotacje do dodania do konta serwisowego.
- **backend.serviceAccount.name**: Nazwa konta serwisowego do użycia. Jeśli nie ustawiono i `create` jest ustawione na true, nazwa jest generowana przy użyciu szablonu fullname.
- **backend.service.type**: Typ usługi do utworzenia dla backendu. Opcje: `ClusterIP`, `NodePort`, `LoadBalancer`.
- **backend.service.port**: Port, na którym usługa backendu będzie dostępna.
- **backend.livenessProbe.httpGet.path**: Ścieżka HTTP używana w liveness probe.
- **backend.livenessProbe.httpGet.port**: Port używany w liveness probe.
- **backend.readinessProbe.httpGet.path**: Ścieżka HTTP używana w readiness probe.
- **backend.readinessProbe.httpGet.port**: Port używany w readiness probe.
- **backend.autoscaling.enabled**: Czy włączyć autoskalowanie dla backendu.
- **backend.autoscaling.minReplicas**: Minimalna liczba replik dla backendu.
- **backend.autoscaling.maxReplicas**: Maksymalna liczba replik dla backendu.
- **backend.autoscaling.targetCPUUtilizationPercentage**: Docelowe zużycie CPU dla autoskalowania.
- **backend.autoscaling.targetMemoryUtilizationPercentage**: Docelowe zużycie pamięci dla autoskalowania.

## Konfiguracja Frontend

- **frontend.replicaCount**: Liczba replik dla wdrożenia frontendu.
- **frontend.image.repository**: Repozytorium obrazu Docker dla frontendu.
- **frontend.image.pullPolicy**: Polityka pobierania obrazu dla frontendu. Opcje: `Always`, `IfNotPresent`, `Never`.
- **frontend.image.tag**: Tag obrazu Docker używanego dla frontendu. Domyślnie: wersja aplikacji w chart.
- **frontend.service.type**: Typ usługi do utworzenia dla frontendu. Opcje: `ClusterIP`, `NodePort`, `LoadBalancer`.
- **frontend.service.port**: Port, na którym usługa frontendu będzie dostępna.
- **frontend.autoscaling.enabled**: Czy włączyć autoskalowanie dla frontendu.
- **frontend.autoscaling.minReplicas**: Minimalna liczba replik dla frontendu.
- **frontend.autoscaling.maxReplicas**: Maksymalna liczba replik dla frontendu.
- **frontend.autoscaling.targetCPUUtilizationPercentage**: Docelowe zużycie CPU dla autoskalowania.
- **frontend.autoscaling.targetMemoryUtilizationPercentage**: Docelowe zużycie pamięci dla autoskalowania.

## Konfiguracja Ingress

- **ingress.enabled**: Czy włączyć ingress dla aplikacji.
- **ingress.className**: Nazwa klasy kontrolera ingress do użycia.
- **ingress.annotations**: Adnotacje do dodania do zasobu ingress.
- **ingress.hosts**: Lista hostów dla zasobu ingress. Każdy host może mieć własny zestaw ścieżek.
- **ingress.tls**: Lista konfiguracji TLS dla zasobu ingress. Każda konfiguracja określa nazwę sekretu i listę hostów.

## Konfiguracja RBAC

- **rbac.create**: Czy utworzyć zasoby RBAC dla aplikacji.
- **rbac.rules**: Lista reguł RBAC do zastosowania. Każda reguła określa grupy API, zasoby i działania (verbs).

# Configuration

This document provides an overview of the configurable values in the [values.yaml](../charts/kam/values.yaml) file for the Helm chart. These values can be customized to configure the deployment of the application. For details how to deploy app refer to [DEPLOYMENT.md](./DEPLOYMENT.md).

## General Configuration

- **nameOverride**: Overrides the name of the chart. This is useful if you want to deploy multiple instances of the chart with different names.
- **fullnameOverride**: Completely overrides the generated name for the chart. This is useful for ensuring consistent naming conventions.

## Global Configuration

### **global.env.KEYCLOAK_JWKS_URL**
- **Description**: The URL for the Keycloak JSON Web Key Set (JWKS). This URL is used to retrieve the public keys for verifying JWT tokens issued by Keycloak.
- **Required**: If `VITE_KEYCLOAK_URL`, `VITE_KEYCLOAK_REALMNAME`, and `VITE_KEYCLOAK_CLIENTNAME` are not provided.
- **Default**: Built from `VITE_KEYCLOAK_URL`, `VITE_KEYCLOAK_REALMNAME`, and `VITE_KEYCLOAK_CLIENTNAME`.
- **Used By**: Backend
- **Example**: `https://keycloak.example.com/realms/myrealm/protocol/openid-connect/certs`

### **global.env.VITE_KEYCLOAK_URL**
- **Description**: The base URL for the Keycloak server. This URL is used by backend to interact with Keycloak for authentication and authorization purposes. It should be accesible from the backend.
- **Required**: If `KEYCLOAK_JWKS_URL` is not provided.
- **Default**: None
- **Used By**: Backend
- **Example**: `https://keycloak.example.com`

### **global.env.FRONTEND_KEYCLOAK_URL**
- **Description**: The URL for the Keycloak server used by the frontend. This URL is specifically used by the frontend for authentication purposes. It should be publicly accesible address.
- **Required**: Yes
- **Default**: None
- **Used By**: Frontend
- **Example**: `https://keycloak.example.com`

### **global.env.VITE_KEYCLOAK_REALMNAME**
- **Description**: The name of the Keycloak realm. A realm in Keycloak is a space where you manage objects such as users, roles, and clients.
- **Required**: Yes
- **Default**: None
- **Used By**: Both Frontend and Backend
- **Example**: `myrealm`

### **global.env.VITE_KEYCLOAK_CLIENTNAME**
- **Description**: The name of the Keycloak client. A client in Keycloak is an entity that can request Keycloak to authenticate a user. You should create a client in Keycloak for our application.
- **Required**: Yes
- **Default**: None
- **Used By**: Both Frontend and Backend
- **Example**: `myclient`

### **global.env.ROLEMAP_NAMESPACE**
- **Description**: The namespace where the role map is stored. This is used to specify the Kubernetes namespace where the role map ConfigMap is located.
- **Required**: No
- **Default**: `default`
- **Used By**: Backend
- **Example**: `mynamespace`

### **global.env.ROLEMAP_NAME**
- **Description**: The name of the role map. This is used to specify the name of the ConfigMap that contains the role map.
- **Required**: No
- **Default**: `role-map`
- **Used By**: Backend
- **Example**: `myrolemap`

## Backend Configuration

- **backend.replicaCount**: The number of replicas for the backend deployment.
- **backend.image.repository**: The Docker image repository for the backend.
- **backend.image.pullPolicy**: The image pull policy for the backend. Options are `Always`, `IfNotPresent`, and `Never`.
- **backend.image.tag**: The tag of the Docker image to use for the backend. Defaults to the chart appVersion.
- **backend.serviceAccount.create**: Whether to create a new service account for the backend.
- **backend.serviceAccount.automount**: Whether to automatically mount the service account's API credentials.
- **backend.serviceAccount.annotations**: Annotations to add to the service account.
- **backend.serviceAccount.name**: The name of the service account to use. If not set and `create` is true, a name is generated using the fullname template.
- **backend.service.type**: The type of service to create for the backend. Options are `ClusterIP`, `NodePort`, and `LoadBalancer`.
- **backend.service.port**: The port on which the backend service will be exposed.
- **backend.livenessProbe.httpGet.path**: The HTTP path to use for the liveness probe.
- **backend.livenessProbe.httpGet.port**: The port to use for the liveness probe.
- **backend.readinessProbe.httpGet.path**: The HTTP path to use for the readiness probe.
- **backend.readinessProbe.httpGet.port**: The port to use for the readiness probe.
- **backend.autoscaling.enabled**: Whether to enable autoscaling for the backend.
- **backend.autoscaling.minReplicas**: The minimum number of replicas for the backend.
- **backend.autoscaling.maxReplicas**: The maximum number of replicas for the backend.
- **backend.autoscaling.targetCPUUtilizationPercentage**: The target CPU utilization percentage for autoscaling.
- **backend.autoscaling.targetMemoryUtilizationPercentage**: The target memory utilization percentage for autoscaling.

## Frontend Configuration

- **frontend.replicaCount**: The number of replicas for the frontend deployment.
- **frontend.image.repository**: The Docker image repository for the frontend.
- **frontend.image.pullPolicy**: The image pull policy for the frontend. Options are `Always`, `IfNotPresent`, and `Never`.
- **frontend.image.tag**: The tag of the Docker image to use for the frontend. Defaults to the chart appVersion.
- **frontend.service.type**: The type of service to create for the frontend. Options are `ClusterIP`, `NodePort`, and `LoadBalancer`.
- **frontend.service.port**: The port on which the frontend service will be exposed.
- **frontend.autoscaling.enabled**: Whether to enable autoscaling for the frontend.
- **frontend.autoscaling.minReplicas**: The minimum number of replicas for the frontend.
- **frontend.autoscaling.maxReplicas**: The maximum number of replicas for the frontend.
- **frontend.autoscaling.targetCPUUtilizationPercentage**: The target CPU utilization percentage for autoscaling.
- **frontend.autoscaling.targetMemoryUtilizationPercentage**: The target memory utilization percentage for autoscaling.

## Ingress Configuration

- **ingress.enabled**: Whether to enable ingress for the application.
- **ingress.className**: The class name of the ingress controller to use.
- **ingress.annotations**: Annotations to add to the ingress resource.
- **ingress.hosts**: A list of hosts for the ingress resource. Each host can have its own set of paths.
- **ingress.tls**: A list of TLS configurations for the ingress resource. Each configuration specifies a secret name and a list of hosts.

## RBAC Configuration

- **rbac.create**: Whether to create RBAC resources for the application.
- **rbac.rules**: A list of RBAC rules to apply. Each rule specifies API groups, resources, and verbs.