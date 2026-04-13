All code should refer to the technical design in the [technical requirement document](https://github.com/sanaul03/ai-sdlc/tree/main/docs/trd).

The code standard:
- Use the Go programming language
- Use PostgreSQL
- No ORM (Object Relational Mapping) framework
- Use REST API as a communication protocol
- Use Golang-migrate and Golang-migrate for pgx5 for any table DDL or data setup, put in a separate SQL file. The general rule is one table per SQL file (including its indices or constraints)

Working process:
- Ensure to use a clean code approach
- Always create unit tests whenever possible
- Ensure the unit test passes before submitting a pull request
- Use conventional commit on pull request title
