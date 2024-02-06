## Answers on questions
1. Used standard library for XML parsing. Suitable for most use cases. In the current case, no additional functionality was required.

If needed more performance or/and xsd validation I would rather use library https://github.com/lestrrat-go/libxml2ы
2. SQL schema create in db/migrations. Use standard data types only.
3. Write an SQL query that finds cycles in a given graph, according to the data model you proposed on item (3).
repository/postgres/graph.go
```
WITH RECURSIVE search_graph(previous_node, next_node, id, depth, path, cycle)
                   AS (
        SELECT e.previous_node, e.next_node, e.id, 1,
               ARRAY[e.id]::varchar[],
               false
        FROM edges e
        UNION ALL
        SELECT e.previous_node, e.next_node, e.id, sg.depth + 1,
               path || e.id,
               e.id = ANY(path)
        FROM edges e, search_graph sg
        WHERE e.previous_node = sg.next_node AND NOT cycle
          and depth < 100
    )
SELECT path
FROM search_graph
WHERE cycle
limit 1;
```
5. Used standard JSON parsing library, encoding/json. Since there are no performance requirements, the standard library is suitable for most use cases. Ease of use and predictability.


## Build and launch service
1. Set DB connection to ENV variables in startup.sh

        `#Local
         export DB_HOST=localhost
         export DB_PORT=5432
         export DB_USER=graph_db_user
         export DB_PASSWORD=graph_db_user
         export DB_NAME=graph
         export DB_SCHEMA=graph
         export SSL_MODE=false`

2. Fill graph.xml
     ```
    <graph>
        <id>g0</id>
        <name>The Graph Name</name>
        <nodes>
            ...
            <node>
                <id>a</id>
                <name>A name</name>
            </node>
            ...
        </nodes>
        <edges>
            ...
            <node>
                <id>a1</id>
                <from>a</from>
                <to>e</to>
                <cost>42</cost>
            </node>
            ...
        </edges>
    </graph>```

3. run startup.sh

    `$sh startup.sh`

STD input Example:
```
{
    "queries": [
        {
            "paths": {
                "start": "a",
                "end": "e"
            }
        },
        {
            "paths": {
                "start": "a",
                "end": "f"
            }
        },
         {
            "paths": {
                "start": "a",
                "end": "z"
            }
        },
        {
            "paths": {
                "start": "a",
                "end": "d"
            }
        },
        {
            "cheapest": {
                "start": "a",
                "end": "d"
            }
        },
        {
            "cheapest": {
                "start": "a",
                "end": "z"
            }
        }
    ]
}
```


