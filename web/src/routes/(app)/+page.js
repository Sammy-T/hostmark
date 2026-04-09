import mdHello from '$lib/hello.md?raw';
import { marked } from 'marked';

// TODO: TEMP
export async function load({ params }) {
    const html = await marked.parse(mdHello);
    return { 
        content: {
            markdown: mdHello,
            html
        }
    };
}
