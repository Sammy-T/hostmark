import mdHello from '$lib/hello.md?raw';
import { marked } from 'marked';

export async function load({ params }) {
    const content = await marked.parse(mdHello);
    return { content };
}
