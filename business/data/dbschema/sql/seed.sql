-- Users
INSERT INTO users ( name, email, roles, password_hash, active, date_created, date_updated) VALUES
	('Admin', 'admin@example.com', '{ADMIN,USER}', '$2a$10$5kQBvJHiwS8sopPUK.fU0eQvKIZ8Ujq9W4roJUZrb6o.51qHc4nVm', true, '2023-06-10 00:00:00', '2023-06-10 00:00:00'),
	('User1', 'user1@example.com', '{USER}', '$2a$10$Hho6jR7msYL70M1n9BH9buFXnqbbT9eqJcBGIGGN9v7Nrz64iNkuO', true, '2023-06-10 00:00:00', '2023-06-10 00:00:00'),
	('User2', 'user2@example.com', '{USER}', '$2a$10$.Kjcriakh9BjVtejapYweuEFPeWs5mTBg4p8kCEd/fOEzZPoboLWK', true, '2023-06-10 00:00:00', '2023-06-10 00:00:00')
    ON CONFLICT DO NOTHING;

-- Permissions
INSERT INTO permissions (name, description)
VALUES ('view_forms', 'Allows users to view forms.'),
       ('create_forms', 'Allows users to create forms.'),
       ('edit_forms', 'Allows users to edit forms.'),
       ('delete_forms', 'Allows users to delete forms.'),
       ('view_responses', 'Allows users to view responses.'),
       ('edit_responses', 'Allows users to edit responses.');

-- User Permissions
INSERT INTO user_permissions (user_id, permission_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (1, 4),
       (1, 5),
       (1, 6),
       (2, 5),
       (2, 6),
       (3, 1),
       (3, 2),
       (3, 3),
       (3, 4);

-- Forms
INSERT INTO forms (form_title, form_description)
VALUES ('Customer Satisfaction Survey', 'A survey to gauge customer satisfaction.'),
       ('Employee Feedback Survey', 'A survey to collect feedback from employees.');

-- Questions
INSERT INTO questions (form_id, question_type, question_text)
VALUES (1, 'text', 'What did you like about our service?'),
       (1, 'text', 'What did you dislike about our service?'),
       (1, 'mcq', 'How would you rate our service?'),
       (2, 'text', 'What do you like about working here?'),
       (2, 'text', 'What do you dislike about working here?'),
       (2, 'mcq', 'How satisfied are you with your job?');

-- Options
INSERT INTO options (question_id, option_text)
VALUES (3, 'Excellent'),
       (3, 'Good'),
       (3, 'Average'),
       (3, 'Poor'),
       (3, 'Terrible'),
       (6, 'Very satisfied'),
       (6, 'Satisfied'),
       (6, 'Neutral'),
       (6, 'Dissatisfied'),
       (6, 'Very dissatisfied');

-- Responses
INSERT INTO responses (form_id, respondent_id)
VALUES (1, 'johndoe'),
       (1, 'janesmith'),
       (2, 'bobjohnson');

-- Answers
INSERT INTO answers (question_id, response_id, answer_text, answer_option_id)
VALUES (1, 1, 'The service was quick and efficient.', NULL),
       (2, 1, 'Nothing, everything was great!', NULL),
       (3, 1, NULL, 1),
       (4, 2, 'I enjoy the team atmosphere.', NULL),
       (5, 2, 'I sometimes find it difficult to balance work and personal life.', NULL),
       (6, 2, NULL, 2),
       (3, 3, NULL, 4),
       (6, 3, NULL, 3);
