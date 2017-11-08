<html>
  <head>
    <title>haha</title>
  </head>
  <body>
    <form action="/login" method="POST" >
      username: <input type="text" name="username"><br/>
      age: <input type="text" name="age" value="1"><br/>
      <input type="submit">
    </form>

    <form action="/upload" method="POST" enctype="multipart/form-data">
      <input type="hidden" name="data" value="123"><br/>
      upload:<input type="file" name="file" /><br/>
      <input type="submit">
    </form>
  </body>
</html>