describe('Articles Page', () => {
  it('should display the title "Latest Articles"', () => {
    cy.visit('http://localhost:3000');

    // Check for the presence of an h2 element with text "Latest Articles"
    cy.contains('h2', 'Latest Articles').should('be.visible');
  });

  it('should display a list of articles', () => {
    cy.visit('http://localhost:3000');

    // Check if the list of articles is present
    cy.get('ul').should('exist');

    // Check if list items are present and at least one list item contains article content
    cy.get('ul li').should('have.length.at.least', 1)
  });
});
