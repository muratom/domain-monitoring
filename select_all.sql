SELECT d.fqdn,
    d.update_at,
    d.update_delay,
    ip4.ip,
    ip6.ip,
    cn.canonical_name,
    mx.host, mx.pref,
    ns.name_server,
    srv.target,
    srv.port,
    srv.priority,
    srv.weight,
    txt.text,
    r.created,
    r.paid_till
FROM domains d
    INNER JOIN ipv4_addresses ip4 ON ip4.domain_id = d.id
    INNER JOIN ipv6_addresses ip6 ON ip6.domain_id = d.id
    INNER JOIN canonical_names cn ON cn.domain_id = d.id
    INNER JOIN mail_exchangers mx ON mx.domain_id = d.id
    INNER JOIN name_servers ns ON ns.domain_id = d.id
    INNER JOIN server_selections srv ON srv.domain_id = d.id
    INNER JOIN text_strings txt ON txt.domain_id = d.id
    INNER JOIN registrations r ON r.domain_id = d.id
WHERE d.fqdn = 'tutu.ru.'