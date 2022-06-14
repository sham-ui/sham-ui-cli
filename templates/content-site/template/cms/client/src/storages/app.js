import createStorage from 'sham-ui-data-storage';

export const { storage, useStorage } = createStorage( {
    routerResolved: false,
    tokenLoaded: false
}, {
    DI: 'app:storage',
    sync: true
} );
