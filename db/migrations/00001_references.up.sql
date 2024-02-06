CREATE TABLE IF NOT EXISTS graphs
(
    id   VARCHAR(64) PRIMARY KEY,
    name VARCHAR(256)
);

comment on column graphs.id is 'graph ID';
comment on column graphs.name is 'graph name';

CREATE TABLE IF NOT EXISTS nodes
(
    id   VARCHAR(64) PRIMARY KEY,
    name VARCHAR(256),
    graph_id varchar(64) references graphs
);

comment on column nodes.id is 'node ID';
comment on column nodes.name is 'node name';


CREATE TABLE IF NOT EXISTS edges
(
    id   VARCHAR(64) PRIMARY KEY,
    previous_node VARCHAR(64) references nodes NOT NULL,
    next_node   VARCHAR(64) references nodes NOT NULL,
    cost NUMERIC(10, 2) default 0

    CONSTRAINT check_previous_not_next CHECK ((previous_node <> next_node))
);

CREATE index if not exists idx_edges_previous_node on edges(previous_node);


comment on column edges.id is 'edge ID';
comment on column edges.previous_node is 'previous_node - from';
comment on column edges.next_node is 'next_node - to';
comment on column edges.cost is 'edge cost';

