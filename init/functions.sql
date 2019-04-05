CREATE FUNCTION user_update(in_nickname VARCHAR(32), in_email VARCHAR(255), in_fullname TEXT, in_about TEXT)
  RETURNS users
  LANGUAGE plpgsql
AS
$BODY$
DECLARE
  user_data users;
BEGIN
  SELECT * INTO STRICT user_data FROM users WHERE LOWER(nickname) = LOWER(in_nickname);

  IF in_email != '' THEN
    user_data.email = in_email;
  END IF;
  IF in_fullname != '' THEN
    user_data.fullname = in_fullname;
  END IF;
  IF in_about != '' THEN
    user_data.about = in_about;
  END IF;

  UPDATE users
  SET email    = user_data.email,
      fullname = user_data.fullname,
      about    = user_data.about
  WHERE id = user_data.id;

  RETURN user_data;
END
$BODY$;
