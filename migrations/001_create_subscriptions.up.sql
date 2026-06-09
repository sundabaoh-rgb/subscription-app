Create table subscriptions (
    id UUID primary key DEFAULT gen_random_uuid(),
    service_name text NOT NULL CHECK (length(trim(service_name)) > 0),
    price integer NOT NULL check (price > 0),
    user_id UUID NOT NULL,
    start_date DATE  NOT NULL,
    end_date DATE  CHECK (end_date IS NULL or end_date > start_date),
    created_at DATE  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_subscription_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscription_end_date ON subscriptions(end_date) where end_date IS NOT NULL;
CREATE INDEX idx_subscription_name ON subscriptions(service_name);