CREATE TABLE IF NOT EXISTS subscriptions (
  id SERIAL PRIMARY KEY,
  user_id int REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  service_name VARCHAR(255) NOT NULL,
  price int NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
