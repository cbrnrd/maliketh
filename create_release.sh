#!/bin/bash

# This script is used to create a new release of the project.

# Things that need to be updated are:
# 1. The version number in server/pyproject.toml
# 2. The version number in client/pyproject.toml
# 3. The CHANGELOG.md file
# 4. The tag

function get_changelog() {
  local currentTag previousTag prevChangelogContents
  currentTag=$(git describe --abbrev=0 --tags "$(git describe --abbrev=0)"^)
  previousTag=$(git describe --abbrev=0)
  prevChangelogContents=$(cat ./CHANGELOG.md)

  {
    echo "## $currentTag";
    echo "";
    git log --pretty=oneline --pretty=format:'`%h` - %s' --abbrev-commit "$currentTag...$previousTag" | grep -v "Upgrade";
    echo "";
  } > CHANGELOG.md
  echo "$prevChangelogContents" >> CHANGELOG.md
}


# Print usage if -h or no args
if [[ $1 == "-h" || $# -eq 0 ]]; then
    echo "Usage: ./create_release.sh <new_version>"
    echo "Do not include the 'v' in the version number!"
    exit 0
fi

# Check if `gh` is installed
if ! command -v gh &> /dev/null; then
    echo "Error: gh is not installed. Please install it."
    exit 1
fi

# Check if `poetry` is installed
if ! command -v poetry &> /dev/null; then
    echo "Error: poetry is not installed. Please install it."
    exit 1
fi

NEW_VERSION=$1

# Update the version number in server/pyproject.toml
cd server && poetry version $NEW_VERSION && cd ..

# Update the version number in client/pyproject.toml
cd client && poetry version $NEW_VERSION && cd ..

# Commit the changes
git add server/pyproject.toml client/pyproject.toml
git commit -m "Bump version to v$NEW_VERSION"

# Create a new tag
git tag -a v$NEW_VERSION

# Push the changes
git push
git push origin v$NEW_VERSION

# Prompt for changelog autogeneration
read -p "Would you like to autogenerate the changelog? (y/n) " AUTOGEN

if [[ $AUTOGEN == "y" ]]; then
    # Generate the changelog
    get_changelog
else
    echo "Write the changelog below (markdown format). Press Ctrl+D when finished."
    prevChangelogContents=$(cat ./CHANGELOG.md)
    cat > CHANGELOG.md
    echo "$prevChangelogContents" >> CHANGELOG.md
fi

read -p "Would you like to add a custom release title? (y/n) " CUSTOM_TITLE

if [[ $CUSTOM_TITLE == "y" ]]; then
    read -p "Enter the release title: " TITLE
    gh release create v$NEW_VERSION --notes-from-tag --title "$TITLE"
else
    gh release create v$NEW_VERSION --notes-from-tag
fi

echo "Done!"