/**
 * Enable animation in app
 * On page load animation disabled for prevent animation on DOM render
 */
export default function enableAnimation() {
    document.querySelector( 'body' ).classList.remove( 'animation-stopped' );
}
