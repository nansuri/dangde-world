# DangDe! World
Dangde World is a web platform for kids where they can play and learn!

Requirements:
- This app is basically Educational app, aiming for Pre-school and School children to learn through play
- User Role:
  - Admin
  - Parents
  - Kids
- This app can have A curriculum which managed by Admin
- Admin will be able to add new learning, curriculum
  - Admin can add new "Activity" which can be under some Categories
  - Categories can be Learning, Playing
  - Sub-Categories can be Learning > Math, Alphabets, etc
- Parents can assign available "Activity"

UI:
- Kahoot like style
- Colorfull, simple, intuitive
- Start with Login page
- Once logged in, depend on the role
  - Kids
    - They will see a sidebar with colorful, icon using some animation
    - They can choose for Math, Alphabets, etc
  - Parents
    - Parent can see statistic of their kids

Original Activity:
- Simple Alphabets
  - Kid learn on how to Say Alphabets
  - Kid can learn how to speak words, available language in English and Bahasa Indonesia
    - i.e. Ba Ca = Baca
- Simple Math 
  - Addition, etc
- Basic Arab Reading
  - Learn to read Basic Arabic character
- Basic Clock reading
  - How to read A Clock, etc

Tech Stack:
- FE: Vue JS
- BE: Go + Gorm + Postgres + Gin Gonic