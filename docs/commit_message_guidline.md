# Commit Message Guidelines

## Format

<type>(<scope>): <subject>

<body>

<footer>

## Examples

feat(category): add permission checks for CRUD operations

- Allow both admin and regular users to manage categories
- Implement authentication middleware for category routes
- Pass requester information to service layer

Closes #123

## Types

- feat: A new feature
- fix: A bug fix
- docs: Documentation only changes
- style: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- refactor: A code change that neither fixes a bug nor adds a feature
- perf: A code change that improves performance
- test: Adding missing tests or correcting existing tests
- chore: Changes to the build process or auxiliary tools and libraries

## Scope

The scope should be the name of the module affected (category, user, restaurant, etc.)

## Subject

- Use the imperative, present tense: "change" not "changed" nor "changes"
- Don't capitalize the first letter
- No dot (.) at the end

## Body

- Use the imperative, present tense
- Include motivation for the change and contrast with previous behavior

## Footer

- Reference issues that this commit closes
