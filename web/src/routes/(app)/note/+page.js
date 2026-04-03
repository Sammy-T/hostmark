import { error } from '@sveltejs/kit';

export const ssr = false

export async function load({ fetch, params }) {
    await fetch('/api/auth/refresh');

    let resp = await fetch(`/api/note/list`);
    if(!resp.ok) error(resp.status, resp.statusText);

    const data = await resp.json();

    return { notes: data };
}
