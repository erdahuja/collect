-- Version: 1.01
-- Description: Create table users table to store user data
CREATE TABLE users (
	user_id       SERIAL   PRIMARY KEY,
	name          TEXT        NOT NULL,
	email         TEXT UNIQUE NOT NULL,
	roles         TEXT[]      NOT NULL,
	password_hash TEXT        NOT NULL,
  active       BOOLEAN     NOT NULL,
	date_created  TIMESTAMP   NOT NULL,
	date_updated  TIMESTAMP   NOT NULL
);


-- Version: 1.02
-- Description: Create table permissions table to store permission data for forms access
CREATE TABLE permissions (
  permission_id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Version: 1.03
-- Description: Create table user_permissions the junction table for the many-to-many relationship between users and permissions
CREATE TABLE user_permissions (
  user_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,
  PRIMARY KEY (user_id, permission_id),
  FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
  FOREIGN KEY (permission_id) REFERENCES permissions (permission_id) ON DELETE CASCADE
);

-- Version: 1.04
-- Description: Create table forms
CREATE TABLE forms (
  form_id SERIAL PRIMARY KEY,
  form_title TEXT NOT NULL,
  form_description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Version: 1.05
-- Description: Create table permission_forms for the many-to-many relationship between permissions and forms
CREATE TABLE permission_forms (
  permission_id INTEGER NOT NULL,
  form_id INTEGER NOT NULL,
  PRIMARY KEY (permission_id, form_id),
  FOREIGN KEY (permission_id) REFERENCES permissions (permission_id) ON DELETE CASCADE,
  FOREIGN KEY (form_id) REFERENCES forms (form_id) ON DELETE CASCADE
);

-- Version: 1.06
-- Description: Create table questions
CREATE TABLE questions (
  question_id SERIAL PRIMARY KEY,
  form_id INTEGER REFERENCES forms(form_id) ON DELETE CASCADE,
  question_type TEXT NOT NULL,
  question_text TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Version: 1.07
-- Description: Create table options
CREATE TABLE options (
  option_id SERIAL PRIMARY KEY,
  question_id INTEGER REFERENCES questions(question_id) ON DELETE CASCADE,
  option_text TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Version: 1.08
-- Description: Create table responses
CREATE TABLE responses (
  response_id SERIAL PRIMARY KEY,
  form_id INTEGER REFERENCES forms(form_id) ON DELETE CASCADE,
  respondent_id TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Version: 1.09
-- Description: Create table answers
CREATE TABLE answers (
  answer_id SERIAL PRIMARY KEY,
  question_id INTEGER REFERENCES questions(question_id) ON DELETE CASCADE,
  response_id INTEGER REFERENCES responses(response_id) ON DELETE CASCADE,
  answer_text TEXT,
  answer_option_id INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Version: 1.10
-- Description: composite index form_respondent_idx
CREATE INDEX form_respondent_idx ON responses (form_id, respondent_id);
