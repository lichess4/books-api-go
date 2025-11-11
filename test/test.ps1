Invoke-WebRequest -Uri "http://localhost:8000/books" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{"title": "Cien años de soledad", "author": "Gabriel García Márquez"}'


Invoke-WebRequest -Uri "http://localhost:8000/books" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{"title": "El Principito", "author": "Antoine de Saint-Exupery"}'


$response = Invoke-WebRequest -Uri "http://localhost:8000/books" -Method GET
$response.Content | ConvertFrom-Json


$response = Invoke-WebRequest -Uri "http://localhost:8000/books/2" -Method GET
$response.Content | ConvertFrom-Json


Invoke-WebRequest -Uri "http://localhost:8000/books/2" `
  -Method PUT `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{"title": "El Principito", "author": "Antoine de Saint-Exupery"}'

$response = Invoke-WebRequest -Uri "http://localhost:8000/books/2" -Method DELETE
$response.Content | ConvertFrom-Json