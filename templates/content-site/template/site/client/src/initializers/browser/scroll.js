/**
 * Process page scroll events
 */
export default function onScroll() {

    // Registry callback
    window.addEventListener( 'scroll', handlers );

    // Call handlers on page load
    handlers();
}

function handlers() {
    navbarShadow();
    backToTop();
}

/**
 * Add shadow to navigation on scroll
 */
function navbarShadow() {
    const navbar = document.querySelector( '.navbar' );
    const needAddShadow = (
        navbar.getBoundingClientRect().top + navbar.ownerDocument.defaultView.pageYOffset > 50
    );
    if ( needAddShadow ) {
        navbar.classList.add( 'navbar-scroll' );
    } else {
        navbar.classList.remove( 'navbar-scroll' );
    }
}

/**
 * Toggle back to top button
 */
function backToTop() {
    const goTopBtn = document.querySelector(  '.back-top' );
    if ( goTopBtn.ownerDocument.defaultView.pageYOffset > 80 ) {
        goTopBtn.classList.add( 'show' );
    } else  {
        goTopBtn.classList.remove( 'show' );
    }
}
