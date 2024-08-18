# Escenario 2

## Requerimientos

- Clúster de Kubernetes (preferentemente con un LoadBalancer configurado).
- Helm v3 instalado en tu máquina local.
- Un motor de contenedores instalado en tu máquina local (opcional).

## Contexto

Una startup tecnológica que desarrolla soluciones de gestión de inventario para pequeñas y medianas empresas ha creado una nueva aplicación en Golang para administrar artículos en un almacén.

La aplicación permite a los usuarios insertar un nuevo artículo, listar todos los artículos existentes y eliminar un artículo específico. Para almacenar los datos, la aplicación utiliza MongoDB como base de datos.

El equipo de DevOps necesita desplegar esta aplicación en su clúster de Kubernetes, junto con MongoDB, utilizando Helm para facilitar la gestión de configuraciones y versiones.

## Solución

### 1. Generar la imagen de la aplicación en Go (opcional)

1. Clona el repositorio con el siguiente comando:

   ```sh
   git clone git@github.com:en-medina/kubernetes-tutorial.git
   ```

2. Accede al directorio que contiene el [código fuente](./src).
3. Inicia sesión en GitHub Container Registry con el siguiente comando:

   ```sh
   docker login --username <username> ghcr.io
   ```

   Necesitarás un GitHub PAT (Personal Access Token). Para más detalles, [consulta este tutorial](https://codefresh.io/docs/docs/integrations/docker-registries/github-container-registry/).

4. Genera la imagen del contenedor y súbela al registro de contenedores:

   ```sh
   # Configura el nombre de usuario de GitHub
   export GH_USERNAME="<username>"

   # Inicia sesión en el registro de contenedores de GitHub
   docker login --username $GH_USERNAME ghcr.io
   [Pega tu token de GitHub en el prompt]

   # Construye y sube la imagen al registro
   docker build -t ghcr.io/$GH_USERNAME/go-mongo-http-basic:0.1.0 ./
   docker push ghcr.io/$GH_USERNAME/go-mongo-http-basic:0.1.0
   ```

5. Ten en cuenta que esta imagen será privada por defecto. Consulta [este tutorial para hacerla pública](https://docs.github.com/en/packages/learn-github-packages/configuring-a-packages-access-control-and-visibility#configuring-access-to-packages-for-your-personal-account).

### 2. Desplegar la aplicación con Helm

1. Accede al directorio que contiene los [archivos de Helm](./helm).
2. Revisa el archivo [values.yaml](./helm/values.yaml) y asegúrate de que tanto la imagen como las credenciales de la base de datos estén configuradas correctamente.
3. Ejecuta los siguientes comandos para generar las dependencias del chart y desplegar la aplicación:

   ```sh
   helm dependency build
   helm install go-app ./
   ```

### 3. Probar la aplicación

1. Localiza la IP o URL desde donde puedes acceder a la aplicación. Si estás usando Minikube, puedes crear un túnel hacia el servicio con:

   ```sh
   minikube service go-app --url
   ```

   Alternativamente, puedes habilitar un port-forward hacia el servicio con:

   ```sh
   kubectl port-forward service/go-app 28015:27017
   ```

   (Asegúrate de validar el puerto donde se ha desplegado).

2. Con el dominio del servicio, puedes probar la funcionalidad de la aplicación con los siguientes comandos:

   ```sh
   # Inserta artículos
   export MY_DOMAIN="127.0.0.1:54639"

   curl -X POST -d '{"name": "car"}' http://$MY_DOMAIN/item
   curl -X POST -d '{"name": "house"}' http://$MY_DOMAIN/item
   curl -X POST -d '{"name": "chair"}' http://$MY_DOMAIN/item

   # Lista todos los artículos
   curl -X GET http://$MY_DOMAIN/items

   # Elimina un artículo
   curl -X DELETE http://$MY_DOMAIN/item/66c170d69c00601f42589626
   ```
