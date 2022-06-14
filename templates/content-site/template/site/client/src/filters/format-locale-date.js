export default date => ( new Date( date ) ).toLocaleDateString( 'ru-RU',  {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
} );
