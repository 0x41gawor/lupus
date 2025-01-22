# Rozdziały
1. Wstęp
    - 1.1 opis z jakimi problemami musi się zmagać operator sieci telco (https://www.ericsson.com/en/blog/2019/4/what-is-closed-loop-automation)
      - sekcja, która: przedstawi dziedzinę (Telekomunikacja), przedstawi motywację do podjęcia tematu (aktualność zagadnienia, znaczenie, rozówj technologiczny, konkretny problem), wskaże lukę badawczą lub problem (brak narzędzia do modelowania i uruchamiania pętli pokazanych przez ENI)
      - sekcja, która: jasno określi co praca ma osiągnąć (jaki jest cel), znaczenie pracy, jakie problemy rozwiązuje, kto może z niej skorzystać
      - sekcja, która: przedstawi zakres, czego ona dotyczy a czego nie
   - 1.2 Opis następnych sekcji i załączników
2. **Stan wiedzy**
   - 2.1 background skąd się wzięła praca	czyli odniesienia do ENI
     - krótki opis ENI
     - krótki opis pętli przedstawionych w [Overview of Prominent Control Loop Architectures](https://www.etsi.org/deliver/etsi_gr/ENI/001_099/017/02.01.01_60/gr_ENI017v020101p.pdf)

2. **Architektura**
   - 2.1 Wstęp
       - 2.1 Opis zamkniętych pętli sterowania znanych z automatyki i robotyki
       - 2.1 Wymagania i Założenia na platformę
   - 2.2 Krótki wstęp czym jest w odniesieniu do pojęć z 2.1 oraz, że w oparciu o Kubernetes
     - Tak jak tu [readme#intro](https://github.com/0x41gawor/lupus?tab=readme-ov-file#lupus)
   - 2.3 Architektura w odniesieniu do 2.1 
     - Tak jak tu [top-level-arch.png](https://github.com/0x41gawor/lupus/blob/master/_img/readme/1.png)
   - 2.4 Podstawowe pojęcia i zasady 
     - Tak jak tu: [detailed-docs](https://github.com/0x41gawor/lupus/blob/master/docs/detailed-docs.md)
     - Odniesienia do [definitions](https://github.com/0x41gawor/lupus/blob/master/docs/defs.md), które będą w załączniku
     - Odniesienia do pełnych specyfikacji [specs/](https://github.com/0x41gawor/lupus/tree/master/docs/spec), które będą w załącznikach
   - 2.5 Krótka instrukcja jak używać
     - Może opis jak tu: [Getting-Started](https://github.com/0x41gawor/lupus/blob/master/docs/getting-started.md)
     - A może same potrzebne kroki jak tu: [readme#how-to-use-it](https://github.com/0x41gawor/lupus/tree/master?tab=readme-ov-file#how-to-use-it)
     > Architektura wg. definicji ENI -> "set of rules and methods that describe the functionality, organization, and implementation of a system"
3. **Implementacja**
   - 3.1 Wstęp
   - 3.2 Mechanizmy stojące za Lupus
     - Opis [controller](https://kubernetes.io/docs/concepts/architecture/controller/) w Kubernetes na buil-in resources
     - Opis [CRD](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/) i [OperatorPattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
     - Opis platformy [Kubebuilder](https://book.kubebuilder.io)
       - Koniecznie architektura oraz flow controllera (operatora)
   - 3.3 Decyzje podjęte podczas developmentu (każda wynikająca z wymagań następnie opis implementacji)
     - To tu objawia się badawcza natura pracy
     - [Komunikacja pomiędzy lupus-elements](https://github.com/0x41gawor/lupus/blob/master/docs/com-bet-lup-ele.md)
     - [Data](https://github.com/0x41gawor/lupus/blob/master/docs/spec/data.md)
     - [Actions](https://github.com/0x41gawor/lupus/blob/master/docs/spec/actions.md)
     - [Polimorfizm w LupN](https://github.com/0x41gawor/lupus/blob/master/docs/go-style-polymorphism.md)
     - [2 workflows](https://github.com/0x41gawor/lupus/blob/master/docs/2-workflows.md)
     - [Open Policy Agents](https://github.com/0x41gawor/lupus/blob/master/docs/open-policy-agents.md)
     - [User-functions](https://github.com/0x41gawor/lupus/blob/master/docs/user-functions.md)
     - Może historia zmian w developmencie? Jak to się zmieniało na przestrzeni sprintów. Tak jak tu: [journal.md](https://github.com/0x41gawor/lupus/blob/master/journal.md)
4. **Test platformy na Open5GS**
    - 4.1 Wstęp
    - 4.2 Opis podobny jak w [examples/open5gs](https://github.com/0x41gawor/lupus/tree/master/examples/open5gs), jednocześnie stanowiący walk-through guide dla usera na żywym przykładzie, krok po kroku, cały kod Ingress-Agent, LupN, Egress-Agent, rysunek workflow i przepływ Date'y w załącznikach.
5. **Wnioski**
   - Co się udało a co nie, czy nadaje się to do użycia itp.
   - Co można poprawić
   - Dalszy rozwój

**Załącznik 1 - definicje**

- [defs.md](https://github.com/0x41gawor/lupus/blob/master/docs/defs.md)

**Załącznik 2 - specyfikacja interfejsów Lupin i Lupout**

- [lupin-lupout.md](https://github.com/0x41gawor/lupus/blob/master/docs/spec/lupin-lupout.md)

**Załącznik 3 - specyfikacja notacji LupN**

- [LupN](https://github.com/0x41gawor/lupus/blob/master/docs/spec/lupn.md)

**Załącznik 4 - Open5GS**

- Wszelki kody źródłowe i rysunki związane z testem na Open5GS


---

Może jakiś state of the art?
- https://www.ericsson.com/en/blog/2019/3/dynamic-orchestration-5g-and-ai-powered-self-healing-networks
- https://www.ericsson.com/en/blog/2019/4/what-is-closed-loop-automation





# Wzorcowa Praca Magisterska wg Chat GPT
Techniczna praca magisterska powinna zawierać konkretne elementy, które spełniają wymagania akademickie oraz pokazują Twoją wiedzę i umiejętności w wybranej dziedzinie. Oto kluczowe sekcje, które powinna zawierać:
---

### **1. Strona tytułowa**
- Tytuł pracy (precyzyjny i odzwierciedlający temat pracy).
- Imię i nazwisko autora.
- Nazwa uczelni, wydziału i kierunku studiów.
- Imię i nazwisko promotora.
- Miejsce i rok.

---

### **2. Streszczenie (Abstrakt)**
- Krótki opis celu pracy, zastosowanej metodyki, wyników i wniosków.
- Powinno być zwięzłe (ok. 200-300 słów).

---

### **3. Spis treści**
- Struktura dokumentu z podziałem na rozdziały i podrozdziały z numeracją stron.

---

### **4. Wstęp**
- Wprowadzenie w tematykę pracy.
- Cel pracy i jej znaczenie.
- Zakres pracy oraz problematyka.
- Krótkie omówienie struktury dokumentu.

---

### **5. Przegląd literatury / Stan wiedzy**
- Przegląd istniejącej literatury, badań i rozwiązań związanych z tematem pracy.
- Analiza dostępnych technologii, metod lub teorii.
- Zidentyfikowanie luki, którą Twoja praca wypełnia.

---

### **6. Metodologia**
- Opis metod, technik i narzędzi użytych w pracy.
- Dlaczego wybrałeś/aś konkretne podejście?
- Specyfikacja założeń, które zostały przyjęte.

---

### **7. Projekt i implementacja**
- Szczegóły techniczne dotyczące rozwiązania opracowanego w pracy.
- Architektura systemu, użyte technologie i narzędzia.
- Diagramy (np. UML, architektoniczne).
- Proces implementacji (np. kod, algorytmy).

---

### **8. Eksperymenty i wyniki**
- Opis przeprowadzonych testów, eksperymentów lub symulacji.
- Wyniki w formie tabel, wykresów i ilustracji.
- Analiza wyników w odniesieniu do założeń i celów pracy.

---

### **9. Dyskusja**
- Krytyczna analiza wyników.
- Porównanie uzyskanych rezultatów z literaturą lub istniejącymi rozwiązaniami.
- Możliwe ograniczenia Twojego podejścia.

---

### **10. Wnioski**
- Podsumowanie pracy.
- Główne osiągnięcia i odpowiedź na postawione pytania badawcze.
- Możliwości rozwoju i przyszłych badań.

---

### **11. Bibliografia**
- Spis cytowanych źródeł (książki, artykuły naukowe, strony internetowe itp.) zgodnie z przyjętym stylem (np. APA, IEEE).

---

### **12. Załączniki (opcjonalne)**
- Dodatkowe materiały, np. fragmenty kodu, dokumentacja techniczna, szczegółowe dane z eksperymentów.

---

### **13. Oświadczenie**
- Oświadczenie o samodzielności wykonania pracy.

---

### Dodatkowe wskazówki:
- **Język**: Zwięzły, techniczny, bez zbędnych ozdobników.
- **Formatowanie**: Zgodne z wymaganiami uczelni (czcionka, marginesy, interlinie).
- **Wizualizacja danych**: Tabele, wykresy, diagramy pomagają w przejrzystej prezentacji wyników.
- **Konsultacje**: Regularnie omawiaj postępy z promotorem.

Czy chcesz szczegółowo opracować któryś z tych elementów? 😊


# Iria Miguel Gonzalez

1. Introduction
   - 1.1 Objectives
   - 1.2 Document structure
2. Technological Background
   - 2.1 Mobile Technologies
     - 2.1.1 Introduction
     - 2.1.2 1G to 3G evolution
     - 2.1.3 4G
     - 2.1.4 5G
   - 2.2 Network Virtualization
     - 2.2.1 Introduction
     - 2.2.2 SDN technology
     - 2.2.3 NFV technology
     - 2.2.4 Network Slicing
3. State of the art
4. Architecture Components
   - 4.1 Introduction
   - 4.2 Docker
     - 4.2.1 Docker Engine
     - 4.2.2 Other Docker components
   - 4.3 Kubernetes
     - 4.3.1 Kubernetes master
     - 4.3.2 Kubernetes workers
     - 4.3.3 Helm
   - 4.4 OSM
     - 4.4.1 Introduction
     - 4.4.2 Architecture
     - 4.4.3 Functionality
   - 4.5 Whole Architecture Overview
5. Developed Solution
   - 5.1 Docker images
     - 5.1.1 Open5GS
     - 5.1.2 Free5GC
   - 5.2 Kubernetes
   - 5.3 OSM
6. Proof of Concept and Results
   - 6.1 Open5GS and SRS-LTE tests
     - 6.1.1 Overview of the EPC Network
     - 6.1.2 Container Images
   - 6.2 Open5GS and UERANSIM tests
     - 6.2.1 Overview of the 5GC Network
     - 6.2.2 Container Images
     - 6.2.3 Run Open5GS and UERANSIM
   - 6.3 Open5GS Slicing tests
     - 6.3.1 Configuration of different AMBR per slice
     - 6.3.2 Selection of UPF and SMF by slice
7. Challenges, conclusions, and future line