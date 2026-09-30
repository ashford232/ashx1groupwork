
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE NOT NULL,
    photo_url VARCHAR(225),
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    title VARCHAR(255) NOT NULL,

    description TEXT,

    priority VARCHAR(25) DEFAULT 'medium',

    due_at TIMESTAMPTZ,

    user_id UUID NOT NULL,

    is_completed BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_todos_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);





CREATE TABLE IF NOT EXISTS labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(50) NOT NULL,

    user_id UUID NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_labels_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_user_label
        UNIQUE (user_id, name)
);

CREATE TABLE IF NOT EXISTS todo_labels (
    todo_id UUID NOT NULL,

    label_id UUID NOT NULL,

    CONSTRAINT pk_todo_labels
        PRIMARY KEY (todo_id, label_id),

    CONSTRAINT fk_todo_labels_todo
        FOREIGN KEY (todo_id)
        REFERENCES todos(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_todo_labels_label
        FOREIGN KEY (label_id)
        REFERENCES labels(id)
        ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS notes(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,  
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP  ,
 is_deleted BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_notes_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notes_labels (
    note_id UUID NOT NULL,

    label_id UUID NOT NULL,

    CONSTRAINT pk_notes_labels
        PRIMARY KEY (note_id, label_id),

    CONSTRAINT fk_note_labels_note
        FOREIGN KEY (note_id)
        REFERENCES notes(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_notes_labels_label
        FOREIGN KEY (label_id)
        REFERENCES labels(id)
        ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS docs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    object_key VARCHAR(100) NOT NULL,
    mime_type VARCHAR(100),
    size INTEGER,
    user_id UUID NOT NULL,
     is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,  
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP  ,

    CONSTRAINT fk_docs_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS docs_labels (
    doc_id UUID NOT NULL,

    label_id UUID NOT NULL,

    CONSTRAINT pk_docs_labels
        PRIMARY KEY (doc_id, label_id),

    CONSTRAINT fk_doc_labels_doc
        FOREIGN KEY (doc_id)
        REFERENCES docs(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_docs_labels_label
        FOREIGN KEY (label_id)
        REFERENCES labels(id)
        ON DELETE CASCADE
);


