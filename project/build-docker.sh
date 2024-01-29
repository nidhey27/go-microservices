#!/bin/bash

# Check if a path and tag are provided as arguments
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <path> --tag <tag>"
    exit 1
fi

path=$1

# Check if the tag flag is provided
if [ "$2" != "--tag" ]; then
    echo "Usage: $0 <path> --tag <tag>"
    exit 1
fi

tag=$3

# Print a message
echo "**** Building Docker Images & Push it to Docker Hub ****"

# Change directory to the provided path
cd "$path" || exit

# List all the directories in the current path and log it to a file
ls -d */ >lsOutput.txt

# Reading the directory names into an array line by line
arr=()
while IFS= read -r line; do
    arr+=("$line")
done <lsOutput.txt

# Looping over directory names array
for ((i = 0; i < ${#arr[@]}; i++)); do
    # Print directory name
    echo "${arr[i]}"

    # Change to the directory
    cd "${arr[i]}" || continue

    # Look for *.dockerfile in the directory
    dockerfiles=$(find . -name "*.dockerfile")

    # Print the found Dockerfiles
    for dockerfile in $dockerfiles; do
        echo "  Found Dockerfile: $dockerfile"

        # Get the service name by reversing the directory name
        service_name=${arr[i]::-1}
        echo "  Service Name: $service_name"

        # Build and tag the Docker image
        docker build -t "nidhey27/$service_name:$tag" -f "$dockerfile" .

        # Push the Docker image to Docker Hub
        docker push "nidhey27/$service_name:$tag"
    done

    # Move back to the original directory
    cd ..
done

# Cleanup: Remove the temporary file
rm lsOutput.txt