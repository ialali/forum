## Database Structure:

### users : Most important table

Column Headings:

- id (Primary Key)
- email (Unique)
- username
- password

### posts : created by users

Column Headings:

- id (Primary Key)
- category
- title
- content
- user_id (Foreign Key ref users)

### comments : created by users

Column Headings:

- id (Primary Key)
- content
- user_id (Foreign Key ref users)
- post_id (Foreign Key ref posts)

### likes_dislikes : interaction by users

Column Headings:

- id (Primary Key)
- like_status (bool)
- user_id (Foreign Key ref users)
- post_id (Foreign Key ref posts)

### categories : selected by users to create posts

Column Headings:

- id (Primary Key)
- category
