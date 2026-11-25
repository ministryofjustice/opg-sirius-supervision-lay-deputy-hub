describe("Manage deputy details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("updates deputy contact details from the deputy details page", () => {
        cy.visit("/supervision/deputies/lay/2");
        cy.contains("a", "Manage deputy contact details").click();
        cy.location("pathname").should("eq", "/supervision/deputies/lay/2/manage-deputy-details");
        cy.get("h1").should("contain", "Manage deputy contact details");

        cy.get('[name="title"]').should("have.value", "Mr").clear();
        cy.get('[name="title"]').type("Dr");
        cy.get('[name="deputy-first-name"]').should("have.value", "Mike").clear();
        cy.get('[name="deputy-first-name"]').type("Alex");
        cy.get('[name="deputy-middle-name"]').should("have.value", "Jonathon").clear();
        cy.get('[name="deputy-middle-name"]').type("Jordan");
        cy.get('[name="deputy-last-name"]').should("have.value", "Deputy").clear();
        cy.get('[name="deputy-last-name"]').type("Smith");
        cy.get('[name="deputy-other-names"]').should("have.value", "Michael").clear();
        cy.get('[name="deputy-other-names"]').type("Alexander");
        cy.get('[name="deputy-date-of-birth"]').should("have.value", "1980-01-01").clear();
        cy.get('[name="deputy-date-of-birth"]').type("1985-02-03");

        cy.get('[name="deputy-mobile-number"]').should("have.value", "0771 2345678").clear();
        cy.get('[name="deputy-mobile-number"]').type("07700 900123");
        cy.get('[name="deputy-phone-number"]').should("have.value", "0115 876 5574").clear();
        cy.get('[name="deputy-phone-number"]').type("020 7946 0123");
        cy.get('[name="deputy-evening-number"]').should("have.value", "0115 2767825").clear();
        cy.get('[name="deputy-evening-number"]').type("020 7946 0456");
        cy.get('[name="deputy-email"]').should("have.value", "testemail@hotmail.co.uk").clear();
        cy.get('[name="deputy-email"]').type("alex@example.com");

        cy.get('[name="address-line-1"]').should("have.value", "Seax House").clear();
        cy.get('[name="address-line-1"]').type("New House");
        cy.get('[name="address-line-2"]').should("have.value", "19 Market Rd").clear();
        cy.get('[name="address-line-2"]').type("20 High Street");
        cy.get('[name="address-line-3"]').should("have.value", "").type("West End");
        cy.get('[name="town"]').should("have.value", "Chelmsford").clear();
        cy.get('[name="town"]').type("London");
        cy.get('[name="county"]').should("have.value", "Essex").clear();
        cy.get('[name="county"]').type("Greater London");
        cy.get('[name="postcode"]').should("have.value", "CM1 1GG").clear();
        cy.get('[name="postcode"]').type("SW1A 1AA");
        cy.get('[name="country"]').should("have.value", "United Kingdom").clear();
        cy.get('[name="country"]').type("France");

        cy.get("#airmail-required-false").should("be.checked");
        cy.get("#airmail-required-true").check();

        cy.contains("button", "Save deputy contact details").click();
        cy.location("pathname").should("eq", "/supervision/deputies/lay/2");
        cy.location("search").should("eq", "?success=deputyDetails");
        cy.get("h1").should("contain", "Deputy details");

        cy.get("#deputy-details").within(() => {
            cy.contains(".govuk-summary-list__key", /^Name$/)
                .siblings(".govuk-summary-list__value").should("contain", "Dr Alex Jordan Smith");
            cy.contains(".govuk-summary-list__key", "Other names")
                .siblings(".govuk-summary-list__value").should("contain", "Alexander");
            cy.contains(".govuk-summary-list__key", "Date of birth")
                .siblings(".govuk-summary-list__value").should("contain", "03/02/1985");
            cy.contains(".govuk-summary-list__key", "Mobile telephone number")
                .siblings(".govuk-summary-list__value").should("contain", "07700 900123");
            cy.contains(".govuk-summary-list__key", "Daytime telephone number")
                .siblings(".govuk-summary-list__value").should("contain", "020 7946 0123");
            cy.contains(".govuk-summary-list__key", "Evening telephone number")
                .siblings(".govuk-summary-list__value").should("contain", "020 7946 0456");
            cy.contains(".govuk-summary-list__key", "Email address")
                .siblings(".govuk-summary-list__value").should("contain", "alex@example.com");
            cy.contains(".govuk-summary-list__key", "Postal address")
                .siblings(".govuk-summary-list__value")
                .should("contain", "New House")
                .and("contain", "20 High Street")
                .and("contain", "West End")
                .and("contain", "London")
                .and("contain", "Greater London")
                .and("contain", "SW1A 1AA");
            cy.contains(".govuk-summary-list__key", "Airmail required")
                .siblings(".govuk-summary-list__value").should("contain", "Yes");
        });
    });

    describe("validation", () => {
        it("displays validation messages", () => {
            cy.visit("/supervision/deputies/lay/2/manage-deputy-details");
            cy.setCookie("fail-route", "contact-details");

            cy.contains("button", "Save deputy contact details").click();

            cy.get(".govuk-error-summary").should("be.visible");
            cy.get(".govuk-error-summary").should("contain", "The building or street must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "Address line 2 must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "Address line 3 must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The country must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The county must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The email must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The deputy first name is required and can't be empty");
            cy.get(".govuk-error-summary").should("contain", "The evening telephone number must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The deputy middle name must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The mobile number must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The deputy other name must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The postcode must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The deputy surname is required and can't be empty");
            cy.get(".govuk-error-summary").should("contain", "The town or city must be 255 characters or fewer");
            cy.get(".govuk-error-summary").should("contain", "The telephone number must be 255 characters or fewer");
        });
    });
});
