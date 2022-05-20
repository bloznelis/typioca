(:require '[clojure.string :as str])

(def exclusions '("“" "”" "gutenberg" "Gutenberg" "Petersburgh" "Mr." "Mrs." "—" "Chapter"))

(defn contains-exclusion? [str]
  (not-every? false? (map #(str/includes? str %) exclusions)))

(defn run [[book-uri output-file]]
  (->> (slurp book-uri)
       (re-seq #"(?s)(.*?(?:\.|\?|!))(?: |$)")
       (map first)
       (filter #(and (> 100 (count %)) (< 15 (count %))))
       (filter  #(not (contains-exclusion? %)))
       (shuffle)
       (map #(str/replace % #"\s+" " "))
       (map str/trim)
       (str/join "\n")
       (#(str/replace % #"’" "'"))
       (#(str/replace % #"_" ""))
       (spit output-file)))

(run *command-line-args*)

(println "Done!")
